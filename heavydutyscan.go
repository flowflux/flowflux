package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
)

// MsgDelimiter ...
var MsgDelimiter = []byte{'\n', '-', '-', '-', '\n'}

// DecodeBase64Message ...
func DecodeBase64Message(encodedMessage []byte) ([]byte, error) {
	decodedMessage := make([]byte, base64.StdEncoding.DecodedLen(len(encodedMessage)))
	_, err := base64.StdEncoding.Decode(decodedMessage, encodedMessage)
	if err != nil {
		return nil, err
	}
	return decodedMessage, nil
}

// HeavyDutyScanner implements the secret sauce that makes it work reliably.
type HeavyDutyScanner struct {
	reader           *bufio.Reader
	delimiter        []byte
	message          []byte
	decodedMessage   []byte
	delimitedMessage []byte
	rest             []byte
	err              error
	Decode           func([]byte) ([]byte, error)
}

// NewHeavyDutyScanner ...
func NewHeavyDutyScanner(reader io.Reader, delimiter []byte) *HeavyDutyScanner {
	return &HeavyDutyScanner{
		reader:    bufio.NewReader(reader),
		delimiter: delimiter,
	}
}

// Scan ...
func (h *HeavyDutyScanner) Scan() bool {
	h.message = nil
	h.decodedMessage = nil
	h.delimitedMessage = nil

	var buffer []byte
	if len(h.rest) > 0 {
		buffer = append(buffer, h.rest...)
		h.rest = nil
	}

	for {
		line, err := h.reader.ReadBytes('\n')
		if err != nil {
			h.err = err
			return false
		}
		buffer = append(buffer, line...)
		idx := bytes.Index(buffer, h.delimiter)
		if idx > -1 {
			h.message = buffer[:idx]
			h.rest = buffer[idx+len(h.delimiter):]
			return true
		}
	}
}

// Message ...
func (h *HeavyDutyScanner) Message() []byte {
	return h.message
}

// DecodedMessage ...
func (h *HeavyDutyScanner) DecodedMessage() ([]byte, error) {
	if h.decodedMessage == nil {
		if h.Decode == nil {
			return nil, fmt.Errorf("No decoder provided")
		}
		decodedMessage, err := h.Decode(h.message)
		if err != nil {
			return nil, err
		}
		h.decodedMessage = decodedMessage
	}
	return h.decodedMessage, nil
}

// DelimitedMessage ...
func (h *HeavyDutyScanner) DelimitedMessage() []byte {
	if h.delimitedMessage == nil {
		h.delimitedMessage = append(h.message, h.delimiter...)
	}
	return h.delimitedMessage
}

// Err ...
func (h *HeavyDutyScanner) Err() error {
	return h.err
}
