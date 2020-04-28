package runner

// InputMessage ...
type InputMessage struct {
	payload []byte
	EOF     bool
}

// collectInputChannels
func collectInputChannels(runners []Runner) []chan<- InputMessage {
	channels := make([]chan<- InputMessage, len(runners))
	for i, r := range runners {
		channels[i] = r.Input()
	}
	return channels
}
