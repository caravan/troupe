package predicate

import "github.com/caravan/troupe/handle"

type (
	// Predicate is a function that returns whether a Message can be handled
	Predicate handle.Handler

	// Accept is a function that accepts a predicated actor.Message
	Accept handle.Handler
)

// Handler constructs a predicated Handler, one where a Predicate check is
// performed on an actor.Message, and if the result is true, will see the
// message passed through the Accept handler
func Handler(p Predicate, a Accept) handle.Handler {
	return handle.And(handle.Handler(p), handle.Handler(a))
}
