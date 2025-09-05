package rlapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// Mock WebSocket server for testing
type MockWSServer struct {
	server       *httptest.Server
	upgrader     websocket.Upgrader
	messages     []string                // Store received messages
	responses    map[string]*PsyResponse // Predefined responses
	pongResponse bool                    // Whether to respond to pings with pong
	dropPongs    bool                    // Whether to drop pong responses
}

func NewMockWSServer() *MockWSServer {
	mock := &MockWSServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		messages:     make([]string, 0),
		responses:    make(map[string]*PsyResponse),
		pongResponse: true, // Default to responding to pings
		dropPongs:    false,
	}

	mock.server = httptest.NewServer(http.HandlerFunc(mock.handleWebSocket))
	return mock
}

func (m *MockWSServer) Close() {
	m.server.Close()
}

func (m *MockWSServer) URL() string {
	return "ws" + strings.TrimPrefix(m.server.URL, "http")
}

func (m *MockWSServer) SetResponse(requestID string, response *PsyResponse) {
	response.ResponseID = requestID
	m.responses[requestID] = response
}

func (m *MockWSServer) SetPongResponse(respond bool) {
	m.pongResponse = respond
}

func (m *MockWSServer) SetDropPongs(drop bool) {
	m.dropPongs = drop
}

func (m *MockWSServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			msg := string(message)
			m.messages = append(m.messages, msg)

			// Handle ping messages
			if strings.Contains(msg, "PsyPing:") && m.pongResponse && !m.dropPongs {
				pongMessage := "PsyPong: \r\n\r\n"
				conn.WriteMessage(websocket.TextMessage, []byte(pongMessage))
			}

			// Parse PsyRequestID from the message
			if strings.Contains(msg, "PsyRequestID:") {
				lines := strings.Split(msg, "\r\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "PsyRequestID:") {
						requestID := strings.TrimSpace(strings.TrimPrefix(line, "PsyRequestID:"))

						// Send predefined response if available
						if response, exists := m.responses[requestID]; exists {
							// Format as PsyNet message with headers and the Result as JSON payload
							psyNetResponse := fmt.Sprintf("PsyTime: %d\r\nPsySig: test_sig\r\nPsyResponseID: %s\r\n\r\n%s",
								time.Now().Unix(), requestID, string(response.Result))
							conn.WriteMessage(websocket.TextMessage, []byte(psyNetResponse))
						}
						break
					}
				}
			}
		}
	}
}

func TestPsyNetRPC_SendRequestSync(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()
	defer rpc.Close()

	// Setup expected response
	expectedResponse := &PsyResponse{
		Result: json.RawMessage(`{"Result":{"shops": [{"id": 1, "name": "test shop"}]}}`),
	}
	mockServer.SetResponse("PsyNetMessage_X_0", expectedResponse)

	// Test sync request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result map[string]interface{}
	err = rpc.sendRequestSync(ctx, "Shops/GetStandardShops", map[string]interface{}{}, &result)
	if err != nil {
		t.Fatalf("sendRequestSync failed: %v", err)
	}

	// Verify result
	if shops, ok := result["shops"]; !ok {
		t.Error("Expected 'shops' in response")
	} else if shopsArray, ok := shops.([]interface{}); !ok || len(shopsArray) != 1 {
		t.Error("Expected shops array with one item")
	}
}

func TestPsyNetRPC_SendRequestAsync(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()
	defer rpc.Close()

	// Setup expected response
	expectedResponse := &PsyResponse{
		Result: json.RawMessage(`{"Result":{"async": "test"}}`),
	}
	mockServer.SetResponse("PsyNetMessage_X_0", expectedResponse)

	// Test async request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	respCh, err := rpc.sendRequestAsync(ctx, "Test/AsyncService", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync failed: %v", err)
	}

	// Wait for response
	var result map[string]interface{}
	err = rpc.awaitResponse(ctx, respCh, &result)
	if err != nil {
		t.Fatalf("awaitResponse failed: %v", err)
	}

	// Verify result
	if async, ok := result["async"]; !ok || async != "test" {
		t.Error("Expected 'async': 'test' in response")
	}
}

func TestPsyNetRPC_ConcurrentRequests(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()
	defer rpc.Close()

	// Setup responses for multiple requests
	for i := 1; i <= 3; i++ {
		response := &PsyResponse{
			Result: json.RawMessage(fmt.Sprintf(`{"Result":{"request": %d}}`, i)),
		}
		mockServer.SetResponse(fmt.Sprintf("PsyNetMessage_X_%d", i-1), response)
	}

	// Send multiple concurrent requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	respCh1, err := rpc.sendRequestAsync(ctx, "Test/Service1", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 1 failed: %v", err)
	}

	respCh2, err := rpc.sendRequestAsync(ctx, "Test/Service2", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 2 failed: %v", err)
	}

	respCh3, err := rpc.sendRequestAsync(ctx, "Test/Service3", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 3 failed: %v", err)
	}

	// Await all responses
	var result1, result2, result3 map[string]interface{}

	if err := rpc.awaitResponse(ctx, respCh1, &result1); err != nil {
		t.Fatalf("awaitResponse 1 failed: %v", err)
	}
	if err := rpc.awaitResponse(ctx, respCh2, &result2); err != nil {
		t.Fatalf("awaitResponse 2 failed: %v", err)
	}
	if err := rpc.awaitResponse(ctx, respCh3, &result3); err != nil {
		t.Fatalf("awaitResponse 3 failed: %v", err)
	}

	// Verify responses match expected values
	if result1["request"] != float64(1) {
		t.Errorf("Expected request 1, got %v", result1["request"])
	}
	if result2["request"] != float64(2) {
		t.Errorf("Expected request 2, got %v", result2["request"])
	}
	if result3["request"] != float64(3) {
		t.Errorf("Expected request 3, got %v", result3["request"])
	}
}

func TestPsyNetRPC_FireAndForgetNoLeak(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()
	defer rpc.Close()

	// Check initial pendingReqs count
	rpc.mu.Lock()
	initialCount := len(rpc.pendingReqs)
	rpc.mu.Unlock()

	// Send requests but never await response (fire-and-forget)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	// Fire multiple requests
	for i := 0; i < 5; i++ {
		_, err := rpc.sendRequestAsync(ctx, "Test/FireForget", map[string]interface{}{})
		if err != nil {
			t.Fatalf("sendRequestAsync %d failed: %v", i, err)
		}
	}

	// Check that pendingReqs grew
	rpc.mu.Lock()
	midCount := len(rpc.pendingReqs)
	rpc.mu.Unlock()

	if midCount <= initialCount {
		t.Error("Expected pendingReqs to grow after sending requests")
	}

	// Cancel context (simulates timeout/cancellation)
	cancel()

	// Wait a bit for cleanup goroutines to run
	time.Sleep(200 * time.Millisecond)

	// Check that pendingReqs was cleaned up
	rpc.mu.Lock()
	finalCount := len(rpc.pendingReqs)
	rpc.mu.Unlock()

	if finalCount != initialCount {
		t.Errorf("Memory leak detected: pendingReqs count %d -> %d -> %d, expected to return to %d",
			initialCount, midCount, finalCount, initialCount)
	}

	t.Logf("Fire-and-forget cleanup verified: %d -> %d -> %d requests",
		initialCount, midCount, finalCount)
}

func TestPsyNetRPC_RawMessage(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	defer rpc.Close()

	eventCh := rpc.Events()

	t.Run("unparseable message", func(t *testing.T) {
		malformedMessage := "this is not a valid psynet message"

		// Verify parseMessage fails
		_, err = rpc.parseMessage(malformedMessage)
		if err == nil {
			t.Error("Expected parseMessage to fail on malformed message")
		}

		// Simulate what readMessages would do
		rpc.sendEvent(EventTypeMessage, malformedMessage)

		// Check that we get the event
		select {
		case event := <-eventCh:
			if event.Type != EventTypeMessage {
				t.Errorf("Expected message event, got %d", event.Type)
			}
			if event.Content != malformedMessage {
				t.Errorf("Expected content to match original message")
			}
		case <-time.After(1 * time.Second):
			t.Error("Expected to receive message event")
		}
	})

	t.Run("unsolicited message", func(t *testing.T) {
		unsolicitedMessage := "PsyTime: 123\r\nPsySig: test\r\n\r\n{\"Result\":{\"unsolicited\":\"data\"}}"

		// Simulate receiving the message
		rpc.sendEvent(EventTypeMessage, unsolicitedMessage)

		// Check that we get the event
		select {
		case event := <-eventCh:
			if event.Type != EventTypeMessage {
				t.Errorf("Expected message event, got %d", event.Type)
			}
			if event.Content != unsolicitedMessage {
				t.Errorf("Expected content to match original message")
			}
		case <-time.After(1 * time.Second):
			t.Error("Expected to receive message event")
		}
	})
}

func TestPsyNetRPC_ConcurrentContextCancellation(t *testing.T) {
	// Setup mock server (no responses, requests will hang)
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance and establish connection
	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	rpc.schedulePing()
	defer rpc.Close()

	// Start multiple requests with different contexts
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each request has its own context with different timeout
			ctx, cancel := context.WithTimeout(context.Background(),
				time.Duration(50+id*10)*time.Millisecond)
			defer cancel()

			var result map[string]interface{}
			err := rpc.sendRequestSync(ctx, fmt.Sprintf("Test/Concurrent%d", id),
				map[string]interface{}{}, &result)
			errors <- err
		}(i)
	}

	wg.Wait()
	close(errors)

	// All should have timed out
	timeoutCount := 0
	for err := range errors {
		if err == context.DeadlineExceeded {
			timeoutCount++
		}
	}

	if timeoutCount != 10 {
		t.Errorf("Expected 10 timeout errors, got %d", timeoutCount)
	}

	// Wait for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify no leaks
	rpc.mu.Lock()
	finalCount := len(rpc.pendingReqs)
	rpc.mu.Unlock()

	if finalCount != 0 {
		t.Errorf("Memory leak: %d pending requests remain", finalCount)
	}
}

func TestPsyNetRPC_ParseMessage(t *testing.T) {
	rpc := &PsyNetRPC{}

	t.Run("valid result", func(t *testing.T) {
		input := fmt.Sprintf("PsyTime: %d\r\nPsySig: test_sig\r\nPsyResponseID: %s\r\n\r\n%s", time.Now().Unix(), "PsyNetMessage_X_1", `{"Result":{"Message":"ok"}}`)

		resp, err := rpc.parseMessage(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.ResponseID != "PsyNetMessage_X_1" {
			t.Errorf("ResponseID = %q, want %q", resp.ResponseID, "PsyNetMessage_X_1")
		}

		var msg struct {
			Message string
		}
		err = json.Unmarshal(resp.Result, &msg)
		if err != nil {
			t.Fatalf("json error: %v", err)
		}
		if msg.Message != "ok" {
			t.Errorf("message = %q, want %q", msg.Message, "ok")
		}
	})

	t.Run("error payload", func(t *testing.T) {
		input := fmt.Sprintf("PsyTime: %d\r\nPsySig: test_sig\r\nPsyResponseID: %s\r\n\r\n%s", time.Now().Unix(), "PsyNetMessage_X_1", `{"Error":{"Type":"InvalidParameters","Message":""}}`)

		resp, err := rpc.parseMessage(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.ResponseID != "PsyNetMessage_X_1" {
			t.Errorf("ResponseID = %q, want %q", resp.ResponseID, "PsyNetMessage_X_1")
		}
		if resp.Error.Type != "InvalidParameters" {
			t.Errorf("Error.Message = %q, want %q", resp.Error.Error(), "InvalidParameters")
		}
	})

	t.Run("missing PsyResponseID header", func(t *testing.T) {
		input := fmt.Sprintf("\r\n\r\n%s", `{"Result":{"Message":"ok"}}`)

		resp, err := rpc.parseMessage(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.ResponseID != "" {
			t.Errorf("ResponseID = %q, want empty", resp.ResponseID)
		}

		var msg struct {
			Message string
		}
		err = json.Unmarshal(resp.Result, &msg)
		if err != nil {
			t.Fatalf("json error: %v", err)
		}
		if msg.Message != "ok" {
			t.Errorf("message = %q, want %q", msg.Message, "ok")
		}
	})
}

func TestPsyNetRPC_BuildMessage(t *testing.T) {
	rpc := &PsyNetRPC{}

	t.Run("request with no body", func(t *testing.T) {
		headers := map[string]string{"PsyPing": ""}
		message, err := rpc.buildMessage(headers, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "PsyPing: \r\n\r\n"
		if message != expected {
			t.Errorf("message = %q, want %q", message, expected)
		}
	})

	t.Run("request with body", func(t *testing.T) {
		headers := map[string]string{
			"PsyService":   "Shops/GetStandardShops",
			"PsyRequestID": "PsyNetMessage_X_123",
		}
		requestData := map[string]interface{}{"test": "data"}

		message, err := rpc.buildMessage(headers, requestData)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should contain all headers including auto-generated PsySig
		if !strings.Contains(message, "PsyService: Shops/GetStandardShops\r\n") {
			t.Error("missing PsyService header")
		}
		if !strings.Contains(message, "PsyRequestID: PsyNetMessage_X_123\r\n") {
			t.Error("missing PsyRequestID header")
		}
		if !strings.Contains(message, "PsySig:") {
			t.Error("missing PsySig header")
		}

		// Should contain JSON body
		if !strings.Contains(message, `{"test":"data"}`) {
			t.Error("missing JSON body")
		}

		// Should have proper structure with delimiter
		if !strings.Contains(message, "\r\n\r\n") {
			t.Error("missing header/body delimiter")
		}
	})
}

func TestPsyNetRPC_IsConnected(t *testing.T) {
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	defer rpc.Close()

	// Should be connected initially
	if !rpc.IsConnected() {
		t.Error("Expected connection to be active initially")
	}

	// Close connection
	rpc.Close()

	// Should not be connected after close
	if rpc.IsConnected() {
		t.Error("Expected connection to be inactive after close")
	}
}

func TestPsyNetRPC_PingPongHandling(t *testing.T) {
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	psyNet := NewPsyNet()
	rpc, err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	go rpc.readMessages()
	defer rpc.Close()

	// Wait a bit to ensure connection is established
	time.Sleep(100 * time.Millisecond)

	// Should be connected
	if !rpc.IsConnected() {
		t.Error("Expected connection to be active")
	}

	// Check that pong channel exists and is ready
	if rpc.pongChan == nil {
		t.Error("Expected pong channel to be initialized")
	}
}
