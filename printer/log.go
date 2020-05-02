package printer

import "os"

// LogLn ...
func LogLn(text string) {
	if len(text) > 512 { // 78
		text = text[:512]
	}
	os.Stderr.WriteString(text + "\n")
}

// ErrLn ...
func ErrLn(text string) {
	os.Stderr.WriteString(text + "\n")
}
