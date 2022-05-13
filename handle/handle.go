package handle

import (
	"fmt"

	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/report"
)

// ErrNonErrorPanic is used to format a recover result that does not support
// the standard error interface
const ErrNonErrorPanic = "panic: %v"

// And constructs a Handler that performs the left Handler, and if it returns
// true will return the result of performing the right Handler
func And(l Handler, r Handler) Handler {
	return func(c actor.Context, m actor.Message) bool {
		if !l(c, m) {
			return false
		}
		return r(c, m)
	}
}

// Or constructs a Handler that performs the left Handler, and if it returns
// false will return the result of performing the right Handler
func Or(l Handler, r Handler) Handler {
	return func(c actor.Context, m actor.Message) bool {
		if l(c, m) {
			return true
		}
		return r(c, m)
	}
}

// Any is a Handler composition that will return true the first time any of its
// constituent Handlers returns true
func Any(first Handler, rest ...Handler) Handler {
	if len(rest) > 0 {
		return Or(first, Any(rest[0], rest[1:]...))
	}
	return first
}

// All is a Handler composition that will only return true if all of its
// constituent Handlers returns true
func All(first Handler, rest ...Handler) Handler {
	if len(rest) > 0 {
		return And(first, All(rest[0], rest[1:]...))
	}
	return first
}

// Panic wraps a Handler and will catch any panic value that is recovered in
// executing that Handler. These will be reported as a report.Error
func Panic(h Handler) Handler {
	return func(c actor.Context, m actor.Message) (result bool) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					report.AnError(c, err)
				} else {
					report.AnError(c, fmt.Errorf(ErrNonErrorPanic, rec))
				}
				result = false
			}
		}()
		return h(c, m)
	}
}

// UnhandledMessage wraps a Handler and will report a report.DeadLetter if that
// Handler returns false
func UnhandledMessage(h Handler) Handler {
	return func(c actor.Context, m actor.Message) bool {
		if !h(c, m) {
			report.AnUnhandledMessage(c, m)
		}
		return true
	}
}
