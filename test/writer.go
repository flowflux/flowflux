package test

import (
	"bytes"
	"log"
)

// Writer ...
type Writer struct {
	Lines   chan []byte
	Running bool
}

// NewWriter ...
func NewWriter() *Writer {
	return &Writer{
		Lines:   make(chan []byte),
		Running: true,
	}
}

func (w *Writer) Write(raw []byte) (int, error) {
	if w.Running {
		b := bytes.Trim(raw, "\n\r")
		if len(b) == 0 {
			return len(raw), nil
		}
		l := make([]byte, len(b))
		copy(l, b)
		w.Lines <- l
	} else {
		log.Fatalf("Writer stopped but received data: \"%v\"\n", string(raw))
	}
	return len(raw), nil
}

// Stop ...
func (w *Writer) Stop(doStop bool) {
	if doStop {
		w.Running = false
		close(w.Lines)
	}
}
