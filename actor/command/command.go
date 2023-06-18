package command

type (
	// Command marks a type as being an Actor Command Message
	Command interface {
		command()
	}

	// Stop instructs an Actor to stop processing Messages entire (shutdown)
	// though an Actor may choose to ignore this Command
	Stop struct{}

	// Pause instructs an Actor to pause its processing of incoming Messages
	// though an Actor may choose to ignore this Command
	Pause struct{}

	// Resume instructs an Actor to resume processing of incoming Messages
	// though an Actor may choose to ignore this Command
	Resume struct{}

	// Restart instructs an Actor to restart processing, including the
	// processing of all Messages it has seen since instantiation, though an
	// Actor may choose to ignore this Command
	Restart struct{}
)

func (*Stop) command()    {}
func (*Pause) command()   {}
func (*Resume) command()  {}
func (*Restart) command() {}
