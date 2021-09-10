package route_test

import (
	"testing"
	"time"

	"github.com/caravan/troupe"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/system"
	"github.com/caravan/troupe/route"
	"github.com/stretchr/testify/assert"
)

type (
	routeTest struct {
		system.System
		t        *testing.T
		handlers int
		router   actor.Address
		msgs     [][]actor.Message
		addrs    []actor.Address
	}

	routeFunc func(
		s actor.Spawner, first actor.Address, rest ...actor.Address,
	) actor.Address
)

func makeRouteTest(t *testing.T, handlers int, router routeFunc) *routeTest {
	res := &routeTest{
		System:   troupe.System(system.Config{}),
		t:        t,
		handlers: handlers,
		msgs:     make([][]actor.Message, handlers),
		addrs:    make([]actor.Address, handlers),
	}

	makeAppender := func(idx int) actor.Factory {
		return actor.Singleton(func(c actor.Context) {
			for msg := range c.Receive() {
				s := res.msgs[idx]
				res.msgs[idx] = append(s, msg)
			}
		})
	}

	for i := 0; i < handlers; i++ {
		res.addrs[i] = res.Spawn(makeAppender(i))
	}

	r := router(res, res.addrs[0], res.addrs[1:]...)
	for i := 0; i < handlers; i++ {
		r.Send() <- i
	}
	time.Sleep(1 * time.Millisecond)
	res.Shutdown()
	return res
}

func (r *routeTest) allMessagesAreLength(l int) {
	as := assert.New(r.t)
	for i := 0; i < r.handlers; i++ {
		as.Equal(l, len(r.msgs[i]))
	}
}

func TestRoundRobin(t *testing.T) {
	rt := makeRouteTest(t, 3, route.RoundRobin)
	rt.allMessagesAreLength(1)
}

func TestFanOut(t *testing.T) {
	rt := makeRouteTest(t, 3, route.FanOut)
	rt.allMessagesAreLength(3)
}
