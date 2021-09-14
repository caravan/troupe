package predicate_test

import (
	"testing"

	"github.com/caravan/essentials/message"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/handle/predicate"
	"github.com/stretchr/testify/assert"
)

type testStruct struct{}

func TestPredicated(t *testing.T) {
	as := assert.New(t)

	var res *testStruct
	h := predicate.Handler(
		predicate.IsA((*testStruct)(nil)),
		func(_ actor.Context, m message.Message) bool {
			res = m.(*testStruct)
			return true
		},
	)

	h(nil, "hello")
	as.Nil(res)

	h(nil, &testStruct{})
	as.NotNil(res)
}
