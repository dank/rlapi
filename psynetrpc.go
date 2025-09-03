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
	requestID   *requestIDCounter
	pendingReqs map[string]chan *PsyResponse
	pongChan    chan struct{}
	connected   bool
	mu          sync.Mutex
	logger      *slog.Logger
}

func newPsyNetRPC(wsConn *websocket.Conn, requestID *requestIDCounter, logger *slog.Logger) *PsyNetRPC {
	return &PsyNetRPC{
		wsConn:      wsConn,
		requestID:   requestID,
		pendingReqs: make(map[string]chan *PsyResponse),
		pongChan:    make(chan struct{}, 1),
		connected:   true,
		logger:      logger,
	}
}

func (p *PsyNetRPC) IsConnected() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.connected && p.wsConn != nil
}

func (p *PsyNetRPC) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.wsConn != nil && p.connected {
		p.connected = false

		p.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

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
	p.logger.Debug("starting ping handler")

	for {
		time.Sleep(pingInterval)

		pingMessage, err := p.buildMessage(map[string]string{"PsyPing": ""}, nil)
		if err != nil {
			p.logger.Error("failed to build ping message", slog.Any("err", err))
			p.mu.Lock()
			p.connected = false
			p.mu.Unlock()
			return
		}

		p.mu.Lock()
		if !p.connected || p.wsConn == nil {
			p.logger.Error("connection lost while preparing to ping")
			p.mu.Unlock()
			return
		}
		if err := p.wsConn.WriteMessage(websocket.TextMessage, []byte(pingMessage)); err != nil {
			p.logger.Error("failed to send ping", slog.Any("err", err))
			p.connected = false
			p.mu.Unlock()
			return
		}
		p.mu.Unlock()

		p.logger.Debug("sent ping")

		select {
		case <-p.pongChan:
			p.logger.Debug("received pong")
		case <-time.After(pongTimeout):
			p.logger.Error("pong timeout reached")
			p.mu.Lock()
			p.connected = false
			p.mu.Unlock()
			return
		}
	}
}

func (p *PsyNetRPC) readMessages() {
	defer func() {
		p.mu.Lock()
		p.connected = false
		p.mu.Unlock()
		p.wsConn.Close()
	}()

	for {
		_, message, err := p.wsConn.ReadMessage()
		if err != nil {
			p.logger.Error("failed to read websocket message", slog.Any("err", err))
			break
		}

		if strings.HasPrefix(string(message), "PsyPong:") {
			select {
			case p.pongChan <- struct{}{}:
			default:
			}
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
	if !p.IsConnected() {
		return nil, fmt.Errorf("websocket connection not established")
	}

	requestID := p.requestID.getID()
	p.logger.Debug("sending websocket request", slog.String("requestID", requestID), slog.String("service", service), slog.Any("data", data))

	respCh := make(chan *PsyResponse, 1)

	headers := map[string]string{
		"PsyService":   service,
		"PsyRequestID": requestID,
	}
	message, err := p.buildMessage(headers, data)
	if err != nil {
		return nil, fmt.Errorf("failed to buildm message: %w", err)
	}

	p.mu.Lock()
	if !p.connected || p.wsConn == nil {
		p.mu.Unlock()
		return nil, fmt.Errorf("connection lost while preparing to send")
	}

	p.pendingReqs[requestID] = respCh
	err = p.wsConn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		delete(p.pendingReqs, requestID)
		p.mu.Unlock()
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	p.mu.Unlock()

	go func() {
		<-ctx.Done()
		p.mu.Lock()
		delete(p.pendingReqs, requestID)
		p.mu.Unlock()
		close(respCh)
	}()

	return respCh, nil
}

func (p *PsyNetRPC) awaitResponse(ctx context.Context, respCh <-chan *PsyResponse, result interface{}) error {
	select {
	case response := <-respCh:
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
