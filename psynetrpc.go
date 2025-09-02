package rlapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// PsyNetRPC represents an authenticated WebSocket connection.
type PsyNetRPC struct {
	wsConn      *websocket.Conn
	pendingReqs map[string]chan *PsyResponse
	mu          sync.Mutex
	logger      *slog.Logger
}

func newPsyNetRPC(wsConn *websocket.Conn, logger *slog.Logger) *PsyNetRPC {
	return &PsyNetRPC{
		wsConn:      wsConn,
		pendingReqs: make(map[string]chan *PsyResponse),
		logger:      logger,
	}
}

func (p *PsyNetRPC) Close() error {
	if p.wsConn != nil {
		return p.wsConn.Close()
	}
	return nil
}

func (p *PsyNetRPC) parseMessage(message string) (*PsyResponse, error) {
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

func (p *PsyNetRPC) buildMessage(headers map[string]string, body interface{}) (string, error) {
	var message strings.Builder
	var jsonData []byte

	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal body: %w", err)
		}

		headers["PsySig"] = generatePsySig(jsonData)
	}

	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	message.WriteString("\r\n")
	message.Write(jsonData)

	return message.String(), nil
}

func (p *PsyNetRPC) pingHandler() {
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

func (p *PsyNetRPC) readMessages() {
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

func (p *PsyNetRPC) sendRequestAsync(ctx context.Context, service string, data interface{}) (<-chan *PsyResponse, error) {
	if p.wsConn == nil {
		return nil, fmt.Errorf("websocket connection not established")
	}

	requestID := getRequestID()
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

func (p *PsyNetRPC) awaitResponse(ctx context.Context, respCh <-chan *PsyResponse, result interface{}) error {
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

func (p *PsyNetRPC) sendRequestSync(ctx context.Context, service string, data interface{}, result interface{}) error {
	respCh, err := p.sendRequestAsync(ctx, service, data)
	if err != nil {
		return err
	}

	return p.awaitResponse(ctx, respCh, result)
}
