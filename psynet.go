package rlapi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

const (
	baseURL      = "https://api.rlpp.psynet.gg/rpc"
	rlFeatureSet = "PrimeUpdate55_1"
	rlVersion    = "250811.43331.492665"
	rlBuildId    = "151471783"
	psySigKey    = "c338bd36fb8c42b1a431d30add939fc7"
)

type psyNetError struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
}

func (e psyNetError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type PsyNet struct {
	client      *http.Client
	wsConn      *websocket.Conn
	requestID   int64
	pendingReqs map[string]chan *PsyResponse
	mu          sync.Mutex
	logger      *slog.Logger
}

type PsyRequest struct {
	Service   string      `json:"PsyService"`
	Sig       string      `json:"PsySig"`
	RequestID string      `json:"PsyRequestID"`
	Data      interface{} `json:"-"`
}

type PsyResponse struct {
	ResponseID string          `json:"PsyResponseID"`
	Result     json.RawMessage `json:"Result"`
	Error      *psyNetError    `json:"Error"`
}

func NewPsyNet() *PsyNet {
	return &PsyNet{
		client:      &http.Client{},
		pendingReqs: make(map[string]chan *PsyResponse),
		logger:      slog.Default(),
	}
}

func (p *PsyNet) establishSocket(url string, psyToken string, sessionID string) error {
	p.logger.Debug("establishing websocket connection", slog.String("url", url))

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, http.Header{
		"PsyBuildID":     []string{rlBuildId},
		"User-Agent":     []string{fmt.Sprintf("RL Win/%s gzip", rlVersion)},
		"PsyEnvironment": []string{"Prod"},
		"PsyToken":       []string{psyToken},
		"PsySessionID":   []string{sessionID},
	})
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}

	p.wsConn = conn

	go p.readMessages()

	return nil
}

func (p *PsyNet) getRequestID() string {
	id := atomic.LoadInt64(&p.requestID)
	atomic.AddInt64(&p.requestID, 1)
	return fmt.Sprintf("PsyNetMessage_X_%d", id)
}

func (p *PsyNet) generatePsySig(body []byte) string {
	h := hmac.New(sha256.New, []byte(psySigKey))
	h.Write([]byte("-"))
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (p *PsyNet) readMessages() {
	defer func() {
		p.wsConn.Close()
	}()

	for {
		_, message, err := p.wsConn.ReadMessage()
		if err != nil {
			p.logger.Error("websocket read error", slog.Any("err", err))
			break
		}

		p.logger.Debug("received websocket response", slog.String("message", string(message)))

		var response PsyResponse
		if err := json.Unmarshal(message, &response); err != nil {
			p.logger.Warn("failed to unmarshal websocket message", slog.Any("err", err), slog.String("message", string(message)))
			continue
		}

		p.mu.Lock()
		ch, exists := p.pendingReqs[response.ResponseID]
		p.mu.Unlock()

		if exists {
			ch <- &response
		}
	}
}

func (p *PsyNet) sendRequestAsync(ctx context.Context, service string, data interface{}) (<-chan *PsyResponse, error) {
	if p.wsConn == nil {
		return nil, fmt.Errorf("websocket connection not established")
	}

	requestID := p.getRequestID()
	p.logger.Debug("sending websocket request", slog.String("requestID", requestID), slog.String("service", service), slog.Any("data", data))

	respCh := make(chan *PsyResponse, 1)

	p.mu.Lock()
	p.pendingReqs[requestID] = respCh
	p.mu.Unlock()

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		delete(p.pendingReqs, requestID)
		p.mu.Unlock()
		close(respCh)
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	headers := map[string]string{
		"PsyService":   service,
		"PsySig":       p.generatePsySig(jsonData),
		"PsyRequestID": requestID,
	}

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	}
	message.WriteString("\n")

	message.Write(jsonData)

	p.logger.Debug("sending websocket request", slog.String("requestID", requestID), slog.String("message", message.String()))

	p.mu.Lock()
	err = p.wsConn.WriteMessage(websocket.TextMessage, []byte(message.String()))
	p.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return respCh, nil
}

func (p *PsyNet) awaitResponse(ctx context.Context, respCh <-chan *PsyResponse, result interface{}) error {
	select {
	case response := <-respCh:
		p.mu.Lock()
		delete(p.pendingReqs, response.ResponseID)
		p.mu.Unlock()

		if response.Error != nil {
			return response.Error
		}

		if err := json.Unmarshal(response.Result, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *PsyNet) sendRequestSync(ctx context.Context, service string, data interface{}, result interface{}) error {
	respCh, err := p.sendRequestAsync(ctx, service, data)
	if err != nil {
		return err
	}

	return p.awaitResponse(ctx, respCh, result)
}

func (p *PsyNet) postJSON(path []string, params interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/%s", baseURL, strings.Join(path, "/"))

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	p.logger.Debug("sending http request", slog.String("url", url), slog.String("body", string(body)))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("RL Win/%s gzip (x86_64-pc-win32) curl-7.67.0 Schannel", rlVersion))
	req.Header.Set("PsyBuildID", rlBuildId)
	req.Header.Set("PsyEnvironment", "Prod")
	req.Header.Set("PsyRequestID", p.getRequestID())
	req.Header.Set("PsySig", p.generatePsySig(body))

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	p.logger.Debug("received http response", slog.String("status", resp.Status), slog.String("body", string(respBytes)))

	var wrapper struct {
		Result json.RawMessage `json:"Result"`
		Error  *psyNetError    `json:"Error"`
	}
	if err := json.Unmarshal(respBytes, &wrapper); err != nil {
		return fmt.Errorf("failed to unmarshal wrapper: %w", err)
	}

	if wrapper.Error != nil {
		return wrapper.Error
	}

	if err := json.Unmarshal(wrapper.Result, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return nil
}
