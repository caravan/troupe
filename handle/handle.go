package handle

import (
	"fmt"

	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/report"
)

// ErrNonErrorPanic is used to format a recover result that does not support
// the standard error interface
const ErrNonErrorPanic = "panic: %v"

// Any is a Handler composition that will return true the first time any of
// its constituent Handlers returns true
func Any(first Handler, rest ...Handler) Handler {
	all := append([]Handler{first}, rest...)
	return func(c actor.Context, m actor.Message) bool {
		for _, h := range all {
			if ok := h(c, m); ok {
				return ok
			}
		}
		return false
	}
}

// All is a Handler composition that will only return true if all of its constituent Handlers returns true
func All(first Handler, rest ...Handler) Handler {
	all := append([]Handler{first}, rest...)
	return func(c actor.Context, m actor.Message) bool {
		for _, h := range all {
			if ok := h(c, m); !ok {
				return ok
			}
		}
		return true
	}
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

// UnhandledMessage wraps a Handler and will report a report.DeadLetter if
// that Handler returns false
func UnhandledMessage(h Handler) Handler {
	return func(c actor.Context, m actor.Message) bool {
		if !h(c, m) {
			report.AnUnhandledMessage(c, m)
		}
		return true
	}
}
