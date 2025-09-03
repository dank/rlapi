package rlapi

import (
	"testing"
)

func TestRequestIDIncrementing(t *testing.T) {
	requestID := requestIDCounter{}
	id0 := requestID.getID()
	id1 := requestID.getID()
	id2 := requestID.getID()

	if id0 != "PsyNetMessage_X_0" {
		t.Errorf("Expected PsyNetMessage_X_0, got %s", id0)
	}
	if id1 != "PsyNetMessage_X_1" {
		t.Errorf("Expected PsyNetMessage_X_1, got %s", id1)
	}
	if id2 != "PsyNetMessage_X_2" {
		t.Errorf("Expected PsyNetMessage_X_2, got %s", id2)
	}
}
