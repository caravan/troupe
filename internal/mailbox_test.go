package internal_test

import (
	"testing"
	"time"

	"github.com/caravan/troupe/internal"
	"github.com/stretchr/testify/assert"
)

func TestMailboxBasics(t *testing.T) {
	as := assert.New(t)

	mb := internal.MakeMailbox()
	as.NotNil(mb)

	done := make(chan struct{})
	go func() {
		msg := <-mb.Receive()
		as.Equal(1, msg)

		msg = <-mb.Receive()
		as.Equal(2, msg)

		close(done)
	}()

	go func() {
		time.Sleep(time.Millisecond)
		mb.Send() <- 1
		mb.Send() <- 2
	}()

	<-done
}

func TestMailboxInOut(t *testing.T) {
	as := assert.New(t)
	mb := internal.MakeMailbox()

	mb.Send() <- "hello"
	mb.Send() <- "there"

	greeting, ok := <-mb.Receive()
	as.True(ok)
	as.Equal("hello", greeting)

	greeting, ok = <-mb.Receive()
	as.True(ok)
	as.Equal("there", greeting)

	mb.Close()
}

func TestMailboxImmediateClose(t *testing.T) {
	as := assert.New(t)
	mb := internal.MakeMailbox()
	mb.Close()

	_, ok := <-mb.Receive() // Nothing anyway
	as.False(ok)
}

func TestMailboxPendingClose(t *testing.T) {
	as := assert.New(t)
	mb := internal.MakeMailbox()
	mb.Send() <- "hello"
	mb.Close()

	_, ok := <-mb.Receive() // Short-circuits pending
	as.False(ok)
}
