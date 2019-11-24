package actor

import (
	"log"

	"github.com/jaqmol/approx/config"
)

// Merge ...
type Merge struct {
	Actor
	ident               string
	sendersCount        int
	closesReceivedCount int
}

// NewMerge ...
func NewMerge(inboxSize int, ident string, sendersCount int) *Merge {
	m := &Merge{
		ident:               ident,
		sendersCount:        sendersCount,
		closesReceivedCount: 0,
	}
	m.init(inboxSize)
	return m
}

// NewMergeFromConf ...
func NewMergeFromConf(inboxSize int, conf *config.Merge) *Merge {
	return NewMerge(inboxSize, conf.Ident, conf.Count)
}

// Start ...
func (m *Merge) Start() {
	if len(m.next) != 1 {
		log.Fatalf(
			"Merge \"%v\" is connected to %v next, 1 expected\n",
			m.ident,
			len(m.next),
		)
	}
	go func() {
		next := m.next[0]
		for message := range m.inbox {
			switch message.Type {
			case DataMessage:
				next.Receive(message)
			case CloseInbox:
				m.closesReceivedCount++
				if m.closesReceivedCount == m.sendersCount {
					next.Receive(message)
					close(m.inbox)
				}
			}
		}
	}()
}
