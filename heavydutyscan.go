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
func (m *HeavyDutyScanner) Scan() bool {
	m.message = nil
	m.decodedMessage = nil
	m.delimitedMessage = nil

	var buffer []byte
	if len(m.rest) > 0 {
		buffer = append(buffer, m.rest...)
		m.rest = nil
	}

	for {
		line, err := m.reader.ReadBytes('\n')
		if err != nil {
			m.err = err
			return false
		}
		buffer = append(buffer, line...)
		idx := bytes.Index(buffer, m.delimiter)
		if idx > -1 {
			m.message = buffer[:idx]
			m.rest = buffer[idx+len(m.delimiter):]
			return true
		}
	}
}

// Message ...
func (m *HeavyDutyScanner) Message() []byte {
	return m.message
}

// DecodedMessage ...
func (m *HeavyDutyScanner) DecodedMessage() ([]byte, error) {
	if m.decodedMessage == nil {
		if m.Decode == nil {
			return nil, fmt.Errorf("No decoder provided")
		}
		decodedMessage, err := m.Decode(m.message)
		if err != nil {
			return nil, err
		}
		m.decodedMessage = decodedMessage
	}
	return m.decodedMessage, nil
}

// DelimitedMessage ...
func (m *HeavyDutyScanner) DelimitedMessage() []byte {
	if m.delimitedMessage == nil {
		m.delimitedMessage = append(m.message, m.delimiter...)
	}
	return m.delimitedMessage
}

// Err ...
func (m *HeavyDutyScanner) Err() error {
	return m.err
}
