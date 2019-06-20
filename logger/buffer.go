package logger

import (
	"bytes"
	"sync"
)

// JSONBuffer is a byte buffer
type JSONBuffer struct {
	v  [][]byte
	mu sync.Mutex
}

// NewJSONBuffer creates a new JSONBuffer
func NewJSONBuffer() *JSONBuffer {
	return &JSONBuffer{
		v: make([][]byte, 0),
	}
}

// Write writes a new line
func (j *JSONBuffer) Write(b []byte) (int, error) {
	j.mu.Lock()
	j.v = append(j.v, b)
	j.mu.Unlock()

	return len(b), nil
}

func (j *JSONBuffer) String() string {
	return "[" + string(bytes.Join(j.v, []byte(","))) + "]"
}
