package report

import (
	"github.com/caravan/troupe/actor"
)

type (
	// Report is a Message that travels upward via the supervisor chain
	Report interface {
		report()
	}

	// Error wraps an error that may propagate upward via the supervisor chain
	Error struct {
		Source  actor.Address
		Wrapped error
	}

	// DeadLetter reports an unhandled Message that may propagate upward via
	// the supervisor chain
	DeadLetter struct {
		Source actor.Address
		actor.Message
	}
)

func (*Error) report()      {}
func (*DeadLetter) report() {}

func (e *Error) Error() string {
	return e.Wrapped.Error()
}

func (e *Error) Unwrap() error {
	return e.Wrapped
}

// AnError reports an error to a Context's supervisor chain
func AnError(c actor.Context, err error) {
	c.Supervisor().Send() <- Error{
		Source:  c.Address(),
		Wrapped: err,
	}
}

// AnUnhandledMessage reports an unhandled actor.Message to a Context's
// supervisor chain
func AnUnhandledMessage(c actor.Context, m actor.Message) {
	c.Supervisor().Send() <- DeadLetter{
		Source:  c.Address(),
		Message: m,
	}
}
