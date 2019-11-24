package event

import (
	"bufio"
	"bytes"
	"io"

	"github.com/jaqmol/approx/config"
)

var msgEndLength int

func init() {
	msgEndLength = len(config.EvntEndBytes)
}

// NewScanner ...
func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitFn)
	return scanner
}

func splitFn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		return 0, nil, io.EOF
	}

	msgEndIndex := bytes.Index(data, config.EvntEndBytes)

	if msgEndIndex == -1 {
		return 0, nil, nil
	}

	token = data[:msgEndIndex]
	advance = len(token) + msgEndLength
	return
}

// ScannedBytesCopy ...
func ScannedBytesCopy(raw []byte) []byte {
	source := bytes.Trim(raw, "\x00")
	destination := make([]byte, len(source))
	copy(destination, source)
	return destination
}

// TerminatedBytesCopy ...
func TerminatedBytesCopy(source []byte) []byte {
	destination := make([]byte, len(source)+msgEndLength)
	copy(destination, source)
	return append(destination, config.EvntEndBytes...)
}
