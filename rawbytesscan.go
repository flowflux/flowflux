package main

import (
	"bufio"
	"io"
)

// RawBytesScanner implements the secret sauce that makes it work reliably.
type RawBytesScanner struct {
	reader  *bufio.Reader
	buffer  []byte
	message []byte
	number  int
	err     error
}

// NewRawBytesScanner ...
func NewRawBytesScanner(reader io.Reader) *RawBytesScanner {
	buffReader := bufio.NewReader(reader)
	return &RawBytesScanner{
		reader: buffReader,
		buffer: make([]byte, buffReader.Size()),
	}
}

// Scan ...
func (r *RawBytesScanner) Scan() bool {
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
func (r *RawBytesScanner) Message() []byte {
	if r.message == nil {
		r.message = make([]byte, r.number)
		copy(r.message, r.buffer)
	}
	return r.message
}

// Err ...
func (r *RawBytesScanner) Err() error {
	return r.err
}
