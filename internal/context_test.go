package internal_test

import (
	"testing"
	"time"

	"github.com/caravan/troupe"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/system"
	"github.com/caravan/troupe/internal"
	"github.com/stretchr/testify/assert"
)

func TestContextBasics(t *testing.T) {
	as := assert.New(t)

	root := internal.MakeContext(nil)
	as.Nil(root.Supervisor())
	as.Equal(root, root.Address())

	done := make(chan struct{})
	var cc actor.Context
	child := root.Spawn(actor.Singleton(func(c actor.Context) {
		cc = c
		as.True(root.EqualTo(c.Supervisor()))
		close(done)
	}))

	<-done
	as.Equal(cc, child)
}

func TestContextClose(t *testing.T) {
	as := assert.New(t)
	sys := troupe.System(system.Config{})
	var ctx actor.Context

	a := sys.Spawn(actor.Singleton(func(c actor.Context) {
		ctx = c
		time.Sleep(time.Millisecond)
	}))
	a.Send() <- "hello"

	time.Sleep(20 * time.Millisecond)
	m, ok := <-ctx.Receive()
	as.Nil(m)
	as.False(ok)
}
