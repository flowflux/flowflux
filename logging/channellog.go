package logging

// ChannelLog ...
type ChannelLog struct {
	Log
}

// NewChannelLog ...
func NewChannelLog(channel chan<- []byte) *ChannelLog {
	l := ChannelLog{}
	l.serialize = make(chan []byte)
	l.dispatchLine = makeDispatchLineWithChannel(channel)
	return &l
}

func makeDispatchLineWithChannel(channel chan<- []byte) func([]byte) {
	return func(line []byte) {
		channel <- line
	}
}
