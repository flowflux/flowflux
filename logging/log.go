package logging

import (
	"bytes"
	"io"

	"github.com/jaqmol/approx/event"
)

// Log ...
type Log struct {
	serialize    chan []byte
	readersCount int
	dispatchLine func(line []byte)
}

// Start ...
func (l *Log) Start() {
	for raw := range l.serialize {
		msg := bytes.Trim(raw, "\x00")
		line := append(msg, '\n')
		l.dispatchLine(line)
	}
}

// Add ...
func (l *Log) Add(reader io.Reader) {
	go l.readFrom(reader)
	l.readersCount++
}

func (l *Log) readFrom(r io.Reader) {
	scanner := event.NewScanner(r)
	for scanner.Scan() {
		original := scanner.Bytes()
		toPassOn := make([]byte, len(original))
		copy(toPassOn, original)
		l.serialize <- toPassOn
	}
	l.readersCount--
	if l.readersCount == 0 {
		close(l.serialize)
	}
}
