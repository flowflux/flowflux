package flowscan

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

const prefixSize = 20

// LengthPrefix implements the secret sauce that makes it work reliably.
type LengthPrefix struct {
	reader   *bufio.Reader
	readBuff []byte
	msgCache []byte
	endIndex uint64
	err      error
}

// NewLengthPrefix ...
func NewLengthPrefix(reader io.Reader) *LengthPrefix {
	buffReader := bufio.NewReader(reader)
	return &LengthPrefix{
		reader:   buffReader,
		readBuff: make([]byte, buffReader.Size()),
	}
}

// Scan ...
func (l *LengthPrefix) Scan() bool {
	l.msgCache = l.msgCache[l.endIndex:]
	l.endIndex = l.parseEndIndex()
	if l.complete() {
		return true
	}
	for {
		fill, err := l.reader.Read(l.readBuff)
		if err != nil {
			l.err = err
			return false
		}
		if fill > 0 {
			l.msgCache = append(l.msgCache, l.readBuff[:fill]...)
			// printer.LogLn(fmt.Sprintf("msgCache = %s", string(l.msgCache)))
		}
		if l.endIndex == 0 {
			l.endIndex = l.parseEndIndex()
		}
		if l.complete() {
			return true
		}
	}
}

func (l *LengthPrefix) parseEndIndex() uint64 {
	if len(l.msgCache) > prefixSize {
		lengthStr := string(l.msgCache[:prefixSize])
		length, err := strconv.ParseUint(lengthStr, 10, 64)
		if err != nil {
			log.Fatalf("Error parsing length prefix: %s\n", lengthStr)
		}
		return prefixSize + length
	}
	return 0
}

func (l *LengthPrefix) complete() bool {
	return l.endIndex != 0 && uint64(len(l.msgCache)) >= l.endIndex
}

// Message ...
func (l *LengthPrefix) Message() []byte {
	return l.msgCache[prefixSize:l.endIndex]
}

// PrefixedMessage ...
func (l *LengthPrefix) PrefixedMessage() []byte {
	return l.msgCache[:l.endIndex]
}

// Err ...
func (l *LengthPrefix) Err() error {
	return l.err
}
