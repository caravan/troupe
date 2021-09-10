package lifecycle

type (
	// Lifecycle marks a type as being an Actor Lifecycle Message
	Lifecycle interface {
		lifecycle()
	}

	lifecycle struct{}

	// Starting informs the Actor that it is in the process of starting
	Starting struct{ lifecycle }

	// Started informs the Actor that it has completed the process of starting
	Started struct{ lifecycle }

	// Paused informs the Actor that it has agreed to being paused
	Paused struct{ lifecycle }

	// Resuming informs the Actor that it is in the process of resuming
	Resuming struct{ lifecycle }

	// Restarting informs the Actor that it is in the process of restarting
	Restarting struct{ lifecycle }

	//Stopped informs the Actor that it has agreed to being stopped
	Stopped struct{ lifecycle }
)

func (*lifecycle) lifecycle() {}
