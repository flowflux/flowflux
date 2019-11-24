package actor

import (
	"io"
	"log"
	"time"
)

// Producer ...
type Producer struct {
	next         []Actable
	performSleep func()
}

// NewProducer ...
func NewProducer(inboxSize int) *Producer {
	return NewThrottledProducer(inboxSize, 0)
}

// NewThrottledProducer ...
func NewThrottledProducer(inboxSize, messagesPerSecond int) *Producer {
	p := &Producer{}
	p.init(inboxSize, messagesPerSecond)
	return p
}

func (p *Producer) init(inboxSize, messagesPerSecond int) {
	var performSleep func()
	if messagesPerSecond > 0 {
		duration := time.Second / time.Duration(messagesPerSecond)
		performSleep = func() {
			time.Sleep(duration)
		}
	} else {
		performSleep = func() {}
	}
	p.next = make([]Actable, 0)
	p.performSleep = performSleep
}

// Produce ...
func (p *Producer) Produce(produce func() ([]byte, error)) error {
	var data []byte
	var err error
	for {
		data, err = produce()
		if err != nil {
			break
		}

		msg := NewDataMessage(data)
		for _, na := range p.next {
			na.Receive(msg)
		}

		p.performSleep()
	}

	p.sendCloseMessage()
	if err == io.EOF {
		return nil
	}
	return err
}

// Next ...
func (p *Producer) Next(next ...Actable) {
	for _, na := range next {
		if na != nil {
			p.next = append(p.next, na)
		}
	}
}

// Receive ...
func (p *Producer) Receive(message Message) {
	log.Fatalln("Producer cannot receive messages")
}

func (p *Producer) sendCloseMessage() {
	for _, na := range p.next {
		na.Receive(NewCloseMessage())
	}
}

// Start ...
func (p *Producer) Start() {}
