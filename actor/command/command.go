package command

type (
	// Command marks a type as being an Actor Command Message
	Command interface {
		command()
	}

	command struct{}

	// Stop instructs an Actor to stop processing Messages entire (shutdown)
	// though an Actor may choose to ignore this Command
	Stop struct{ command }

	// Pause instructs an Actor to pause its processing of incoming Messages
	// though an Actor may choose to ignore this Command
	Pause struct{ command }

	// Resume instructs an Actor to resume processing of incoming Messages
	// though an Actor may choose to ignore this Command
	Resume struct{ command }

	// Restart instructs an Actor to restart processing, including the
	// processing of all Messages it has seen since instantiation, though an
	// Actor may choose to ignore this Command
	Restart struct{ command }
)

func (*command) command() {}
