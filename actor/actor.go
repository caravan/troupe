package actor

import "github.com/caravan/essentials/message"

type (
	// Actor is the bare minimum interface for an Actor
	Actor func(Context)

	// Factory is responsible for instantiating new Actors based on any
	// provided set of Arg
	Factory interface {
		New(...Arg) Actor
	}

	// Arg is an argument to the instantiation of a new Actor
	Arg interface{}

	// Spawner is a type that is capable of contextually Spawning a new Actor
	Spawner interface {
		Spawn(Factory, ...Arg) Address
	}

	// Context allows an Actor to contextually interact with the Actor System
	Context interface {
		Spawner
		message.Receiver
		Address() Address
		Supervisor() Address
	}

	// Address is the address of an Actor. This is a location-independent
	// interface that allows a client of an Actor to deliver a Message to it
	Address interface {
		message.Sender
		EqualTo(Address) bool
	}

	// Singleton is an Actor that implements Factory in such a way that it
	// will return itself whenever New is invoked
	Singleton Actor
)

// New allows a Singleton to act as a Factory of itself. Args are ignored
func (s Singleton) New(_ ...Arg) Actor {
	return Actor(s)
}
