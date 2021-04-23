package netrasp

// ConfigResult is returned by a Configure operation and contains output and results.
type ConfigResult struct {
	ConfigCommands []ConfigCommand
}

// ConfigCommand contains a configuration command together with the output from that command.
type ConfigCommand struct {
	// The command that was sent to the device
	Command string
	// The output seen after entering the command
	Output string
}
