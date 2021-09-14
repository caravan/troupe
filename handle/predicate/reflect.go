package predicate

import (
	"reflect"

	"github.com/caravan/essentials/message"
	"github.com/caravan/troupe/actor"
)

// IsA returns a Predicate that compares a provided message.Message with the
// resolved type of the configured interface. This Predicate uses the Go
// reflection facility, and will eventually be replaced by generics
func IsA(i interface{}) Predicate {
	t := reflect.TypeOf(i)
	return func(_ actor.Context, m message.Message) bool {
		return reflect.TypeOf(m) == t
	}
}
