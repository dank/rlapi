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
	server    *httptest.Server
	upgrader  websocket.Upgrader
	messages  []string                // Store received messages
	responses map[string]*PsyResponse // Predefined responses
}

func NewMockWSServer() *MockWSServer {
	mock := &MockWSServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		messages:  make([]string, 0),
		responses: make(map[string]*PsyResponse),
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

func TestPsyNet_SendRequestSync(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance
	psyNet := NewPsyNet()

	// Connect to mock server
	err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	// Setup expected response
	expectedResponse := &PsyResponse{
		Result: json.RawMessage(`{"shops": [{"id": 1, "name": "test shop"}]}`),
	}
	mockServer.SetResponse("PsyNetMessage_X_0", expectedResponse)

	// Test sync request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result map[string]interface{}
	err = psyNet.sendRequestSync(ctx, "Shops/GetStandardShops", map[string]interface{}{}, &result)
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

func TestPsyNet_SendRequestAsync(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance
	psyNet := NewPsyNet()

	// Connect to mock server
	err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	// Setup expected response
	expectedResponse := &PsyResponse{
		Result: json.RawMessage(`{"async": "test"}`),
	}
	mockServer.SetResponse("PsyNetMessage_X_0", expectedResponse)

	// Test async request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	respCh, err := psyNet.sendRequestAsync(ctx, "Test/AsyncService", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync failed: %v", err)
	}

	// Wait for response
	var result map[string]interface{}
	err = psyNet.awaitResponse(ctx, respCh, &result)
	if err != nil {
		t.Fatalf("awaitResponse failed: %v", err)
	}

	// Verify result
	if async, ok := result["async"]; !ok || async != "test" {
		t.Error("Expected 'async': 'test' in response")
	}
}

func TestPsyNet_ConcurrentRequests(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance
	psyNet := NewPsyNet()

	// Connect to mock server
	err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	// Setup responses for multiple requests
	for i := 1; i <= 3; i++ {
		response := &PsyResponse{
			Result: json.RawMessage(fmt.Sprintf(`{"request": %d}`, i)),
		}
		mockServer.SetResponse(fmt.Sprintf("PsyNetMessage_X_%d", i-1), response)
	}

	// Send multiple concurrent requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	respCh1, err := psyNet.sendRequestAsync(ctx, "Test/Service1", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 1 failed: %v", err)
	}

	respCh2, err := psyNet.sendRequestAsync(ctx, "Test/Service2", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 2 failed: %v", err)
	}

	respCh3, err := psyNet.sendRequestAsync(ctx, "Test/Service3", map[string]interface{}{})
	if err != nil {
		t.Fatalf("sendRequestAsync 3 failed: %v", err)
	}

	// Await all responses
	var result1, result2, result3 map[string]interface{}

	if err := psyNet.awaitResponse(ctx, respCh1, &result1); err != nil {
		t.Fatalf("awaitResponse 1 failed: %v", err)
	}
	if err := psyNet.awaitResponse(ctx, respCh2, &result2); err != nil {
		t.Fatalf("awaitResponse 2 failed: %v", err)
	}
	if err := psyNet.awaitResponse(ctx, respCh3, &result3); err != nil {
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

func TestPsyNet_FireAndForgetNoLeak(t *testing.T) {
	// Setup mock server
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance
	psyNet := NewPsyNet()

	// Connect to mock server
	err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

	// Check initial pendingReqs count
	psyNet.mu.Lock()
	initialCount := len(psyNet.pendingReqs)
	psyNet.mu.Unlock()

	// Send requests but never await response (fire-and-forget)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	// Fire multiple requests
	for i := 0; i < 5; i++ {
		_, err := psyNet.sendRequestAsync(ctx, "Test/FireForget", map[string]interface{}{})
		if err != nil {
			t.Fatalf("sendRequestAsync %d failed: %v", i, err)
		}
	}

	// Check that pendingReqs grew
	psyNet.mu.Lock()
	midCount := len(psyNet.pendingReqs)
	psyNet.mu.Unlock()

	if midCount <= initialCount {
		t.Error("Expected pendingReqs to grow after sending requests")
	}

	// Cancel context (simulates timeout/cancellation)
	cancel()

	// Wait a bit for cleanup goroutines to run
	time.Sleep(200 * time.Millisecond)

	// Check that pendingReqs was cleaned up
	psyNet.mu.Lock()
	finalCount := len(psyNet.pendingReqs)
	psyNet.mu.Unlock()

	if finalCount != initialCount {
		t.Errorf("Memory leak detected: pendingReqs count %d -> %d -> %d, expected to return to %d",
			initialCount, midCount, finalCount, initialCount)
	}

	t.Logf("Fire-and-forget cleanup verified: %d -> %d -> %d requests",
		initialCount, midCount, finalCount)
}

func TestPsyNet_RequestIDIncrementing(t *testing.T) {
	psyNet := NewPsyNet()

	// Test that requestID increments properly
	id1 := psyNet.getRequestID()
	id2 := psyNet.getRequestID()
	id3 := psyNet.getRequestID()

	if id1 != "PsyNetMessage_X_0" {
		t.Errorf("Expected first request ID to be 'PsyNetMessage_X_0', got '%s'", id1)
	}
	if id2 != "PsyNetMessage_X_1" {
		t.Errorf("Expected second request ID to be 'PsyNetMessage_X_1', got '%s'", id2)
	}
	if id3 != "PsyNetMessage_X_2" {
		t.Errorf("Expected third request ID to be 'PsyNetMessage_X_2', got '%s'", id3)
	}
}

func TestPsyNet_ConcurrentContextCancellation(t *testing.T) {
	// Setup mock server (no responses, requests will hang)
	mockServer := NewMockWSServer()
	defer mockServer.Close()

	// Create PsyNet instance
	psyNet := NewPsyNet()

	// Connect to mock server
	err := psyNet.establishSocket(mockServer.URL(), "test-token", "test-session")
	if err != nil {
		t.Fatalf("Failed to establish socket: %v", err)
	}

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
			err := psyNet.sendRequestSync(ctx, fmt.Sprintf("Test/Concurrent%d", id),
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
	psyNet.mu.Lock()
	finalCount := len(psyNet.pendingReqs)
	psyNet.mu.Unlock()

	if finalCount != 0 {
		t.Errorf("Memory leak: %d pending requests remain", finalCount)
	}
}
