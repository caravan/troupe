package route

import "github.com/caravan/troupe/actor"

// RoundRobin spawns a new actor.Actor that performs round-robin routing of
// each incoming message.Message
func RoundRobin(
	s actor.Spawner, first actor.Address, rest ...actor.Address,
) actor.Address {
	return s.Spawn(actor.Singleton(func(c actor.Context) {
		addr := append([]actor.Address{first}, rest...)
		idx := 0
		max := len(addr)
		for {
			msg := <-c.Receive()
			addr[idx].Send() <- msg
			idx++
			if idx == max {
				idx = 0
			}

		}
	}))
}

// FanOut spawns a new actor.Actor that performs found-out routing of each
// incoming message.Message
func FanOut(
	s actor.Spawner, first actor.Address, rest ...actor.Address,
) actor.Address {
	return s.Spawn(actor.Singleton(func(c actor.Context) {
		addr := append([]actor.Address{first}, rest...)
		for {
			msg := <-c.Receive()
			for _, a := range addr {
				a.Send() <- msg
			}
		}
	}))
}
