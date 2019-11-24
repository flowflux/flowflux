package logging

import (
	"io"
	"log"
)

// WriterLog ...
type WriterLog struct {
	Log
}

// NewWriterLog ...
func NewWriterLog(w io.Writer) *WriterLog {
	l := WriterLog{}
	l.serialize = make(chan []byte)
	l.dispatchLine = makeDispatchLineWithWriter(w)
	return &l
}

func makeDispatchLineWithWriter(writer io.Writer) func([]byte) {
	return func(line []byte) {
		n, err := writer.Write(line)
		if err != nil {
			log.Fatalln("WriterLog error:", err)
		}
		if n != len(line) {
			panic("Couldn't write complete line")
		}
	}
}
