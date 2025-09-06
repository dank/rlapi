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

type EventType int

const (
	EventTypeDisconnected EventType = iota
	EventTypeMessage
)

// Event represents connection events or raw messages from the server
type Event struct {
	Type    EventType
	Content string
}

// PsyNetRPC represents an authenticated WebSocket connection.
type PsyNetRPC struct {
	wsConn *websocket.Conn
	mu     sync.Mutex
	logger *slog.Logger

	pingTimer   *time.Timer
	pongChan    chan struct{}
	eventCh     chan *Event
	pendingReqs map[string]chan *PsyResponse

	requestID *requestIDCounter
	connected bool
}

func newPsyNetRPC(wsConn *websocket.Conn, requestID *requestIDCounter, logger *slog.Logger) *PsyNetRPC {
	return &PsyNetRPC{
		wsConn:      wsConn,
		requestID:   requestID,
		pendingReqs: make(map[string]chan *PsyResponse),
		pongChan:    make(chan struct{}, 1),
		eventCh:     make(chan *Event, 32),
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
	if p.wsConn != nil && p.connected {
		_ = p.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		err := p.wsConn.Close()

		// clean up pending events
		p.connected = false

		if p.pingTimer != nil {
			p.pingTimer.Stop()
			p.pingTimer = nil
		}

		for reqID, ch := range p.pendingReqs {
			close(ch)
			delete(p.pendingReqs, reqID)
		}

		p.mu.Unlock()
		p.sendEvent(EventTypeDisconnected, "")

		return err
	}

	p.mu.Unlock()
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

func (p *PsyNetRPC) schedulePing() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.connected || p.wsConn == nil {
		return
	}

	p.pingTimer = time.AfterFunc(pingInterval, p.sendPing)
}

func (p *PsyNetRPC) sendPing() {
	pingMessage, err := p.buildMessage(map[string]string{"PsyPing": ""}, nil)
	if err != nil {
		p.logger.Error("failed to build ping message", slog.Any("err", err))
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
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.logger.Debug("sent ping")

	select {
	case <-p.pongChan:
		p.logger.Debug("received pong")
		p.schedulePing()
	case <-time.After(pongTimeout):
		p.logger.Error("pong timeout reached")
		_ = p.Close()
		return
	}
}

func (p *PsyNetRPC) readMessages() {
	defer func() {
		_ = p.Close()
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
			p.sendEvent(EventTypeMessage, string(message))
			continue
		}

		if response.ResponseID != "" {
			p.mu.Lock()
			ch, exists := p.pendingReqs[response.ResponseID]
			p.mu.Unlock()

			if exists {
				ch <- response
				continue
			}
		}

		p.sendEvent(EventTypeMessage, string(message))
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
		ch := p.pendingReqs[requestID]
		delete(p.pendingReqs, requestID)
		p.mu.Unlock()

		if ch != nil {
			close(ch)
		}
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

// Events returns a channel that receives raw messages and connection events
func (p *PsyNetRPC) Events() <-chan *Event {
	return p.eventCh
}

func (p *PsyNetRPC) sendEvent(eventType EventType, content string) {
	select {
	case p.eventCh <- &Event{
		Type:    eventType,
		Content: content,
	}:
	default:
		p.logger.Warn("event channel is full, dropping event",
			slog.String("type", fmt.Sprintf("%d", eventType)),
			slog.String("content", content))
	}
}
