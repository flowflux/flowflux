package processor

import (
	"github.com/flowflux/flowflux/config"
)

var evntEndLength int

func init() {
	evntEndLength = len(config.EvntEndBytes)
}

func evntEndedCopy(data []byte) []byte {
	dataCopy := make([]byte, len(data)+evntEndLength)
	copy(dataCopy, data)
	return append(dataCopy, config.EvntEndBytes...)
}
