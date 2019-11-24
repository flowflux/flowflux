package actor

import (
	"log"

	"github.com/jaqmol/approx/config"
)

// Fork ...
type Fork struct {
	Actor
	ident          string
	receiversCount int
}

// NewFork ...
func NewFork(inboxSize int, ident string, receiversCount int) *Fork {
	f := &Fork{
		ident:          ident,
		receiversCount: receiversCount,
	}
	f.init(inboxSize)
	return f
}

// NewForkFromConf ...
func NewForkFromConf(inboxSize int, conf *config.Fork) *Fork {
	return NewFork(inboxSize, conf.Ident, conf.Count)
}

// Start ...
func (f *Fork) Start() {
	if len(f.next) < 2 {
		log.Fatalf(
			"Fork \"%v\" is connected to %v next only, minimum is 2\n",
			f.ident,
			len(f.next),
		)
	}
	go func() {
		for message := range f.inbox {
			for _, na := range f.next {
				na.Receive(message)
			}
			if message.Type == CloseInbox {
				close(f.inbox)
			}
		}
	}()
}
