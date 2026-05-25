package rlapi

import (
	"sync"
	"testing"
)

func TestRequestID(t *testing.T) {
	counter := requestIDCounter{}

	if id := counter.getID(); id != "PsyNetMessage_X_0" {
		t.Fatalf("expected PsyNetMessage_X_0, got %s", id)
	}
	if id := counter.getID(); id != "PsyNetMessage_X_1" {
		t.Fatalf("expected PsyNetMessage_X_1, got %s", id)
	}

	const n = 50
	results := make([]string, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = counter.getID()
		}()
	}
	wg.Wait()

	seen := make(map[string]bool, n)
	for _, id := range results {
		if seen[id] {
			t.Errorf("duplicate request ID: %s", id)
		}
		seen[id] = true
	}
}
