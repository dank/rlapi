package rlapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	baseURL      = "https://api.rlpp.psynet.gg/rpc"
	rlFeatureSet = "PrimeUpdate55_1"
	rlVersion    = "250811.43331.492665"
	rlBuildId    = "151471783"
)

type psyNetError struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
}

func (e psyNetError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type PsyNet struct {
	client *http.Client
	wsConn *websocket.Conn
}

func NewPsyNet() *PsyNet {
	return &PsyNet{
		client: &http.Client{},
	}
}

func (p *PsyNet) establishSocket(url string, psyToken string, sessionID string) error {
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	headers := http.Header{
		"PsyBuildID":     []string{rlBuildId},
		"User-Agent":     []string{fmt.Sprintf("RL Win/%s gzip", rlVersion)},
		"PsyEnvironment": []string{"Prod"},
		"PsyToken":       []string{psyToken},
		"PsySessionID":   []string{sessionID},
	}

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}

	p.wsConn = conn

	return nil
}

func (p *PsyNet) postJSON(path []string, params interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/%s", baseURL, strings.Join(path, "/"))

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("RL Win/%s gzip (x86_64-pc-win32) curl-7.67.0 Schannel", rlVersion))

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
