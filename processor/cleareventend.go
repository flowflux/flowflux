package processor

import (
	"bytes"

	"github.com/jaqmol/approx/config"
)

// ClearEventEnd ...
func ClearEventEnd(raw []byte) []byte {
	msg := bytes.ReplaceAll(raw, []byte("\x00"), []byte(""))
	msgEndIndex := bytes.Index(msg, config.EvntEndBytes)
	if msgEndIndex == -1 {
		return msg
	}
	return msg[:msgEndIndex]
}
