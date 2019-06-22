package logger

import (
	"sync"
)

// JSONBuffer is a byte buffer
type JSONBuffer struct {
	v  []byte
	mu sync.Mutex
}

// NewJSONBuffer creates a new JSONBuffer
func NewJSONBuffer() *JSONBuffer {
	return &JSONBuffer{
		v: []byte{},
	}
}

// Write writes a new line
func (j *JSONBuffer) Write(b []byte) (int, error) {
	// Trim the new line
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}

	j.mu.Lock()
	j.v = append(j.v, b...)
	j.v = append(j.v, ',')
	j.mu.Unlock()

	return len(b), nil
}

func (j *JSONBuffer) String() string {
	j.mu.Lock()
	defer j.mu.Unlock()

	return "[" + string(j.v[:len(j.v)-1]) + "]"
}
