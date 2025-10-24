package rlapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	baseURL      = "https://api.rlpp.psynet.gg/rpc"
	gameVersion  = "251020.62592.500294"
	featureSet   = "PrimeUpdate56_1"
	psyBuildId   = "717252342"
	psySigKey    = "c338bd36fb8c42b1a431d30add939fc7"
	pingInterval = 20 * time.Second
	pongTimeout  = 10 * time.Second
)

type psyNetError struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
}

func (e psyNetError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// PsyNet represents the HTTP API client, see PsyNetRPC for the WebSocket client.
type PsyNet struct {
	client    *http.Client
	requestID *requestIDCounter
	logger    *slog.Logger
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

func generatePsySig(body []byte) string {
	h := hmac.New(sha256.New, []byte(psySigKey))
	h.Write([]byte("-"))
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func NewPsyNet() *PsyNet {
	return &PsyNet{
		client:    &http.Client{},
		requestID: &requestIDCounter{},
		logger:    slog.Default(),
	}
}

func NewPsyNetWithLogger(logger *slog.Logger) *PsyNet {
	return &PsyNet{
		client:    &http.Client{},
		requestID: &requestIDCounter{},
		logger:    logger,
	}
}

func (p *PsyNet) establishSocket(url string, playerID PlayerID, psyToken string, sessionID string) (*PsyNetRPC, error) {
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
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}

	return newPsyNetRPC(conn, playerID, p.requestID, p.logger), nil
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
	req.Header.Set("PsyRequestID", p.requestID.getID())
	req.Header.Set("PsySig", generatePsySig(body))

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
