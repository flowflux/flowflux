package logging

import "io"

// Logger ...
type Logger interface {
	Start()
	Add(io.Reader)
}
