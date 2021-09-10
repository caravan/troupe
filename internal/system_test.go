package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/caravan/troupe"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/report"
	"github.com/caravan/troupe/actor/system"
	"github.com/caravan/troupe/internal"
	"github.com/stretchr/testify/assert"
)

func TestInternalSystem(t *testing.T) {
	as := assert.New(t)

	sys := internal.MakeSystem(system.Config{})
	as.Nil(sys.Supervisor())
	as.Equal(sys.Context, sys.Address())

	sys.Send() <- report.DeadLetter{
		Source:  nil,
		Message: "ignored",
	}

	done := make(chan struct{})
	var cc actor.Context
	child := sys.Spawn(actor.Singleton(func(c actor.Context) {
		cc = c
		as.True(sys.Address().EqualTo(c.Supervisor()))
		close(done)
	}))

	<-done
	as.Equal(cc, child)
}

func TestActorSystem(t *testing.T) {
	as := assert.New(t)
	s := troupe.System(system.Config{})
	as.NotNil(s)

	var msg actor.Message
	done := make(chan struct{})

	addr := s.Spawn(
		actor.Singleton(func(c actor.Context) {
			msg = <-c.Receive()
			close(done)
		}),
	)

	go func() {
		time.Sleep(time.Millisecond * 10)
		addr.Send() <- "hello"
	}()

	<-done
	as.Equal("hello", msg)
}

func TestSystemRouting(t *testing.T) {
	as := assert.New(t)

	var errors []actor.Message
	var deadLetters []actor.Message

	sys := troupe.System(system.Config{
		Errors: actor.Singleton(func(c actor.Context) {
			for msg := range c.Receive() {
				errors = append(errors, msg)
			}
		}),
		DeadLetters: actor.Singleton(func(c actor.Context) {
			for msg := range c.Receive() {
				deadLetters = append(deadLetters, msg)
			}
		}),
	})
	as.NotNil(sys)

	sys.Spawn(actor.Singleton(func(c actor.Context) {
		c.Supervisor().Send() <- report.Error{
			Source:  c.Address(),
			Wrapped: fmt.Errorf("blew up once"),
		}
		c.Supervisor().Send() <- report.Error{
			Source:  c.Address(),
			Wrapped: fmt.Errorf("blew up twice"),
		}
		c.Supervisor().Send() <- report.DeadLetter{
			Source:  c.Address(),
			Message: "random message I didn't handle",
		}
	}))

	time.Sleep(time.Millisecond)
	as.Equal(2, len(errors))
	as.Equal(1, len(deadLetters))
	sys.Shutdown()
}
