package internal

import (
	"sync"

	"github.com/caravan/troupe/actor"
)

// Context is the internal implementation of both an actor.Context (and an
// actor.Address). It also serves as the foundation for the internal
// implementation of system.System
type Context struct {
	sync.RWMutex
	*Mailbox
	supervisor actor.Address
	children   map[actor.Address]bool
}

// MakeContext instantiates a new Context
func MakeContext(supervisor actor.Address) *Context {
	return &Context{
		Mailbox:    MakeMailbox(),
		supervisor: supervisor,
		children:   map[actor.Address]bool{},
	}
}

// Address returns the actor.Address that is managed by this Context
func (c *Context) Address() actor.Address {
	return c
}

// Supervisor returns the actor.Address of the Supervisor of this Context
func (c *Context) Supervisor() actor.Address {
	return c.supervisor
}

// Spawn instantiates a new actor.Actor as a child of this Context
func (c *Context) Spawn(f actor.Factory, args ...actor.Arg) actor.Address {
	child := MakeContext(c)
	go func() {
		c.register(child)
		f.New(args...)(child)
		c.unregister(child)
		child.Close()
	}()
	return child
}

func (c *Context) register(a actor.Address) {
	c.Lock()
	defer c.Unlock()
	c.children[a] = true
}

func (c *Context) unregister(a actor.Address) {
	c.Lock()
	defer c.Unlock()
	delete(c.children, a)
}

// EqualTo compares the two actor.Address instances for equality
func (c *Context) EqualTo(other actor.Address) bool {
	return c == other
}
