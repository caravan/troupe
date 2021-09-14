package handle

import (
	"github.com/caravan/essentials/message"
	"github.com/caravan/troupe/actor"
)

// Handler is a partial actor.Actor that can be composed with others
type Handler func(actor.Context, message.Message) bool

// New allows a Handler to act as a Factory of itself. Args are ignored and
// both panic and dead letter reporting are performed
func (h Handler) New(_ ...actor.Arg) actor.Actor {
	wrapped := Panic(UnhandledMessage(h))
	return func(c actor.Context) {
		for m := range c.Receive() {
			_ = wrapped(c, m)
		}
	}
}
