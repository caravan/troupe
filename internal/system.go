package internal

import (
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/report"
	"github.com/caravan/troupe/actor/system"
)

// System is the internal implementation of system.System
type System struct {
	*Context
	system.Config
}

var devNull = actor.Singleton(func(c actor.Context) {
	for range c.Receive() {
		// do nothing with it
	}
})

// MakeSystem instantiates a new System
func MakeSystem(cfg system.Config) *System {
	sys := &System{
		Context: MakeContext(nil),
		Config:  cfg,
	}
	go sys.start()
	return sys
}

// Shutdown instructs the System to stop processing Messages and to destroy
// all of its actor.Actor children
func (s *System) Shutdown() {
	s.Close()
}

func (s *System) start() {
	deadLetters := s.spawnOrDevNull(s.DeadLetters)
	errors := s.spawnOrDevNull(s.Errors)

	for msg := range s.Receive() {
		switch msg := msg.(type) {
		case report.DeadLetter:
			deadLetters.Send() <- msg
		case report.Error:
			errors.Send() <- msg
		}
	}
}

func (s *System) spawnOrDevNull(f actor.Factory) actor.Address {
	if f != nil {
		return s.Spawn(f)
	}
	return s.Spawn(devNull)
}
