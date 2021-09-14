package internal

import (
	"container/list"

	"github.com/caravan/essentials/message"
)

// Mailbox is the queueing mechanism for a local actor.Actor
type Mailbox struct {
	in     chan message.Message
	out    chan message.Message
	closed chan struct{}
	queue  *list.List
}

// MakeMailbox instantiates a new Mailbox
func MakeMailbox() *Mailbox {
	res := &Mailbox{
		queue:  list.New(),
		in:     make(chan message.Message),
		out:    make(chan message.Message),
		closed: make(chan struct{}),
	}
	go res.process()
	return res
}

func (m *Mailbox) process() {
	for {
		if elem := m.queue.Front(); elem != nil {
			select {
			case <-m.closed:
				goto closed
			case m.out <- elem.Value:
				m.queue.Remove(elem)
			case msg := <-m.in:
				m.queue.PushBack(msg)
			}
		} else {
			select {
			case <-m.closed:
				goto closed
			case msg := <-m.in:
				m.queue.PushBack(msg)
			}
		}
	}
closed:
	close(m.in)
	close(m.out)
}

// Send returns the Mailbox's sending channel. Will usually be exposed by an
// actor.Address implementation that composes it in
func (m *Mailbox) Send() chan<- message.Message {
	return m.in
}

// Receive returns the Mailbox's receiving channel. Will usually be exposed by
// an actor.Context implementation that composes it in
func (m *Mailbox) Receive() <-chan message.Message {
	return m.out
}

// Close the Mailbox
func (m *Mailbox) Close() {
	close(m.closed)
}
