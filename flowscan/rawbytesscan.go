package flowscan

import (
	"bufio"
	"io"
)

// RawBytes implements the secret sauce that makes it work reliably.
type RawBytes struct {
	reader  *bufio.Reader
	buffer  []byte
	message []byte
	number  int
	err     error
}

// NewRawBytes ...
func NewRawBytes(reader io.Reader) *RawBytes {
	buffReader := bufio.NewReader(reader)
	return &RawBytes{
		reader: buffReader,
		buffer: make([]byte, buffReader.Size()),
	}
}

// Scan ...
func (r *RawBytes) Scan() bool {
	r.message = nil
	for {
		n, err := r.reader.Read(r.buffer)
		r.number = n
		if err != nil {
			r.err = err
			return false
		}
		if r.number > 0 {
			return true
		}
	}
}

// Message ...
func (r *RawBytes) Message() []byte {
	if r.message == nil {
		r.message = make([]byte, r.number)
		copy(r.message, r.buffer)
	}
	return r.message
}

// Err ...
func (r *RawBytes) Err() error {
	return r.err
}
