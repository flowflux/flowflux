package printer

import "os"

// LogLn ...
func LogLn(text string) {
	if len(text) > 78 {
		text = text[:78]
	}
	os.Stderr.WriteString(text + "\n")
}

// ErrLn ...
func ErrLn(text string) {
	os.Stderr.WriteString(text + "\n")
}
