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
	"time"

	"github.com/gorilla/websocket"
)

const (
	baseURL      = "https://api.rlpp.psynet.gg/rpc"
	gameVersion  = "250811.43331.492665"
	featureSet   = "PrimeUpdate55_1"
	psyBuildId   = "151471783"
	psySigKey    = "c338bd36fb8c42b1a431d30add939fc7"
	pingInterval = 20 * time.Second
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
	ctx         context.Context
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

func NewPsyNet(ctx context.Context) *PsyNet {
	return &PsyNet{
		client:      &http.Client{},
		pendingReqs: make(map[string]chan *PsyResponse),
		logger:      slog.Default(),
		ctx:         ctx,
	}
}

func (p *PsyNet) Close() error {
	if p.wsConn != nil {
		return p.wsConn.Close()
	}
	return nil
}

func (p *PsyNet) establishSocket(url string, psyToken string, sessionID string) error {
	p.logger.Debug("establishing websocket connection", slog.String("url", url))

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url, http.Header{
		"PsyBuildID":     []string{psyBuildId},
		"User-Agent":     []string{fmt.Sprintf("RL Win/%s gzip", gameVersion)},
		"PsyEnvironment": []string{"Prod"},
		"PsyToken":       []string{psyToken},
		"PsySessionID":   []string{sessionID},
	})
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}

	p.wsConn = conn

	go p.readMessages()
	go p.pingHandler()

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

func (p *PsyNet) parseMessage(message string) (*PsyResponse, error) {
	delimiter := "\r\n\r\n"
	index := strings.Index(message, delimiter)
	if index == -1 {
		return nil, fmt.Errorf("message does not contain expected delimiter")
	}

	headersPart := message[:index]
	jsonPayload := message[index+len(delimiter):]

	headers := make(map[string]string)
	headerLines := strings.Split(headersPart, "\r\n")
	for _, line := range headerLines {
		if colonIndex := strings.Index(line, ":"); colonIndex != -1 {
			key := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])
			headers[key] = value
		}
	}

	responseID := headers["PsyResponseID"]

	var jsonResult PsyResponse
	if err := json.Unmarshal([]byte(jsonPayload), &jsonResult); err != nil {
		return nil, fmt.Errorf("failed to parse json payload: %w", err)
	}

	jsonResult.ResponseID = responseID

	return &jsonResult, nil
}

func (p *PsyNet) buildMessage(headers map[string]string, body interface{}) (string, error) {
	var message strings.Builder
	var jsonData []byte

	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal body: %w", err)
		}

		headers["PsySig"] = p.generatePsySig(jsonData)
	}

	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	message.WriteString("\r\n")
	message.Write(jsonData)

	return message.String(), nil
}

func (p *PsyNet) pingHandler() {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		if p.wsConn != nil {
			pingMessage, err := p.buildMessage(map[string]string{"PsyPing": ""}, nil)
			if err != nil {
				p.logger.Error("failed to build ping message", slog.Any("err", err))
				p.mu.Unlock()
				return
			}

			if err := p.wsConn.WriteMessage(websocket.TextMessage, []byte(pingMessage)); err != nil {
				p.logger.Error("failed to send psynet ping", slog.Any("err", err))
				p.mu.Unlock()
				return
			}

			p.logger.Debug("sent psynet ping")
		}
		p.mu.Unlock()
	}
}

func (p *PsyNet) readMessages() {
	defer func() {
		p.wsConn.Close()
	}()

	for {
		_, message, err := p.wsConn.ReadMessage()
		if err != nil {
			p.logger.Error("failed to read websocket message", slog.Any("err", err))
			break
		}

		if strings.HasPrefix(string(message), "PsyPong:") {
			p.logger.Debug("received psynet pong")
			continue
		}

		p.logger.Debug("received websocket response", slog.String("message", string(message)))

		response, err := p.parseMessage(string(message))
		if err != nil {
			p.logger.Error("failed to parse psynet message", slog.Any("err", err), slog.String("message", string(message)))
			continue
		}

		if response.ResponseID != "" {
			p.mu.Lock()
			ch, exists := p.pendingReqs[response.ResponseID]
			p.mu.Unlock()

			if exists {
				ch <- response
			}
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

	headers := map[string]string{
		"PsyService":   service,
		"PsyRequestID": requestID,
	}

	message, err := p.buildMessage(headers, data)
	if err != nil {
		return nil, err
	}

	p.logger.Debug("sending websocket request", slog.String("requestID", requestID), slog.String("message", message))

	p.mu.Lock()
	err = p.wsConn.WriteMessage(websocket.TextMessage, []byte(message))
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
	req.Header.Set("User-Agent", fmt.Sprintf("RL Win/%s gzip (x86_64-pc-win32) curl-7.67.0 Schannel", gameVersion))
	req.Header.Set("PsyBuildID", psyBuildId)
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
