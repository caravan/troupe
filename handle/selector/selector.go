package selector

import (
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/handle"
)

// Selector is used to map an incoming actor.Message into a form that can
// be passed downstream into a handle.Handler
type Selector func(actor.Context, actor.Message) actor.Message

// Into returns a handle.Handler that performs a Selector mapping and then
// passes the resulting value into the provided handler
func Into(s Selector, h handle.Handler) handle.Handler {
	return func(c actor.Context, m actor.Message) bool {
		return h(c, s(c, m))
	}
}
