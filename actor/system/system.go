package system

import "github.com/caravan/troupe/actor"

type (
	// System is the root container for all actor.Actor instances
	System interface {
		actor.Address
		actor.Spawner
		Shutdown()
	}

	// Config describes necessary configuration for an Actor System
	Config struct {
		DeadLetters actor.Factory
		Errors      actor.Factory
	}

	// Constructor is the standard interface for instantiating a System
	Constructor func(config Config) System
)
