package handle_test

import (
	"errors"
	"testing"
	"time"

	"github.com/caravan/essentials/message"
	"github.com/caravan/troupe"
	"github.com/caravan/troupe/actor"
	"github.com/caravan/troupe/actor/system"
	"github.com/caravan/troupe/handle"
	"github.com/stretchr/testify/assert"
)

type handlerTest struct {
	sys         system.System
	deadLetters int
	errors      int
}

func makeHandlerTest() *handlerTest {
	res := &handlerTest{}
	res.sys = troupe.System(system.Config{
		DeadLetters: actor.Singleton(func(c actor.Context) {
			for range c.Receive() {
				res.deadLetters++
			}
		}),
		Errors: actor.Singleton(func(c actor.Context) {
			for range c.Receive() {
				res.errors++
			}
		}),
	})
	return res
}

var (
	returnTrue = func(_ actor.Context, _ message.Message) bool {
		return true
	}
	returnFalse = func(_ actor.Context, _ message.Message) bool {
		return false
	}
	returnError = func(_ actor.Context, _ message.Message) bool {
		panic(errors.New("there was an error"))
	}
	returnNonError = func(_ actor.Context, _ message.Message) bool {
		panic("this is a string")
	}
)

func TestAll(t *testing.T) {
	as := assert.New(t)
	test := makeHandlerTest()

	allTrue := test.sys.Spawn(handle.All(returnTrue, returnTrue))
	firstTrue := test.sys.Spawn(handle.All(returnTrue, returnFalse))
	firstFalse := test.sys.Spawn(handle.All(returnFalse, returnTrue))
	errorsOut := test.sys.Spawn(handle.All(returnTrue, returnError))

	allTrue.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(0, test.deadLetters)
	as.Equal(0, test.errors)

	firstTrue.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(1, test.deadLetters)
	as.Equal(0, test.errors)

	firstFalse.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(2, test.deadLetters)
	as.Equal(0, test.errors)

	errorsOut.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(2, test.deadLetters)
	as.Equal(1, test.errors)
}

func TestAny(t *testing.T) {
	as := assert.New(t)
	test := makeHandlerTest()

	firstTrue := test.sys.Spawn(handle.Any(returnTrue, returnFalse))
	secondTrue := test.sys.Spawn(handle.Any(returnFalse, returnTrue))
	noneTrue := test.sys.Spawn(handle.Any(returnFalse, returnFalse))
	errorsOut := test.sys.Spawn(handle.Any(returnFalse, returnError))

	firstTrue.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(0, test.deadLetters)
	as.Equal(0, test.errors)

	secondTrue.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(0, test.deadLetters)
	as.Equal(0, test.errors)

	noneTrue.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(1, test.deadLetters)
	as.Equal(0, test.errors)

	errorsOut.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(1, test.deadLetters)
	as.Equal(1, test.errors)
}

func TestNonErrorPanic(t *testing.T) {
	as := assert.New(t)
	test := makeHandlerTest()

	errorsOut := test.sys.Spawn(handle.All(returnTrue, returnNonError))

	errorsOut.Send() <- "hello"
	time.Sleep(time.Millisecond)
	as.Equal(0, test.deadLetters)
	as.Equal(1, test.errors)
}
