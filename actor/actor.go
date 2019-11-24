package actor

import "log"

// Actor ...
type Actor struct {
	inbox chan Message
	next  []Actable
}

// NewActor ...
func NewActor(inboxSize int) *Actor {
	a := &Actor{}
	a.init(inboxSize)
	return a
}

func (a *Actor) init(inboxSize int) {
	a.inbox = make(chan Message, inboxSize)
	a.next = make([]Actable, 0)
}

// Next ...
func (a *Actor) Next(next ...Actable) {
	for _, na := range next {
		if na != nil {
			a.next = append(a.next, na)
		}
	}
}

// Receive ...
func (a *Actor) Receive(message Message) {
	a.inbox <- message
}

// Start ...
func (a *Actor) Start() {
	log.Fatalln("Actor.Start() must be overwritten")
}
