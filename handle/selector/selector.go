package selector

import (
	"github.com/caravan/essentials/message"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/handle"
)

// Selector is used to map an incoming message.Message into a form that can be
// passed downstream into a handle.Handler
type Selector func(actor.Context, message.Message) message.Message

// Into returns a handle.Handler that performs a Selector mapping and then
// passes the resulting value into the provided handler
func Into(s Selector, h handle.Handler) handle.Handler {
	return func(c actor.Context, m message.Message) bool {
		return h(c, s(c, m))
	}
}
