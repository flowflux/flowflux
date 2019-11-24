package test

import (
	"strings"
	"testing"

	"github.com/jaqmol/approx/event"
)

// CheckCmdLogMsg ...
func CheckCmdLogMsg(t *testing.T, expectedPrefix string, data []byte) {
	msg, err := event.UnmarshalLogMsg(data)
	logMsg, cmdErr, err := msg.PayloadOrError()
	if err != nil {
		t.Fatal(err)
	}
	if logMsg != nil {
		if !strings.HasPrefix(*logMsg, expectedPrefix) {
			t.Fatalf("Unexpected command log message: %v", *logMsg)
		}
	}
	if cmdErr != nil {
		t.Fatal(cmdErr.Error())
	}
}
