package lifecycle

type (
	// Lifecycle marks a type as being an Actor Lifecycle Message
	Lifecycle interface {
		lifecycle()
	}

	// Starting informs the Actor that it is in the process of starting
	Starting struct{}

	// Started informs the Actor that it has completed the process of starting
	Started struct{}

	// Paused informs the Actor that it has agreed to being paused
	Paused struct{}

	// Resuming informs the Actor that it is in the process of resuming
	Resuming struct{}

	// Restarting informs the Actor that it is in the process of restarting
	Restarting struct{}

	//Stopped informs the Actor that it has agreed to being stopped
	Stopped struct{}
)

func (*Starting) lifecycle()   {}
func (*Started) lifecycle()    {}
func (*Paused) lifecycle()     {}
func (*Resuming) lifecycle()   {}
func (*Restarting) lifecycle() {}
func (*Stopped) lifecycle()    {}
