package rlapi

import (
	"fmt"
	"sync/atomic"
)

type requestIDCounter struct {
	value int64
}

func (r *requestIDCounter) getID() string {
	id := atomic.LoadInt64(&r.value)
	atomic.AddInt64(&r.value, 1)
	return fmt.Sprintf("PsyNetMessage_X_%d", id)
}
