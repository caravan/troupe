package selector_test

import (
	"testing"

	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/handle/selector"
	"github.com/stretchr/testify/assert"
)

type selectorTest struct {
	name string
	age  int
}

func TestInto(t *testing.T) {
	as := assert.New(t)

	var n string
	h := selector.Into(
		func(_ actor.Context, m actor.Message) actor.Message {
			return m.(*selectorTest).name
		},
		func(_ actor.Context, m actor.Message) bool {
			n = m.(string)
			return true
		},
	)

	h(nil, &selectorTest{
		name: "bob",
		age:  49,
	})
	as.Equal("bob", n)
}
