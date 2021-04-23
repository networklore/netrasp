package netrasp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// SROS is the Netrasp driver for Nokia SR OS (MD-CLI).
type sros struct {
	Connection connection
}

// Close connection to device.
func (s sros) Close(ctx context.Context) error {
	s.Connection.Close(ctx)

	return nil
}

// Configure device.
func (s sros) Configure(ctx context.Context, commands []string) (ConfigResult, error) {
	var result ConfigResult

	_, err := s.Run(ctx, "edit-config exclusive")
	if err != nil {
		return result, fmt.Errorf("unable to enter exclusive edit mode: %w", err)
	}
	for _, command := range commands {
		output, err := s.Run(ctx, command)
		configCommand := ConfigCommand{Command: command, Output: output}
		result.ConfigCommands = append(result.ConfigCommands, configCommand)
		if err != nil {
			return result, fmt.Errorf("unable to run command '%s': %w", command, err)
		}
	}
	_, err = s.Run(ctx, "commit")
	if err != nil {
		return result, fmt.Errorf("unable to commit configuration: %w", err)
	}

	_, err = s.Run(ctx, "quit-config")
	if err != nil {
		return result, fmt.Errorf("unable to quit from configuration mode: %w", err)
	}

	return result, nil
}

// Dial opens a connection to a device.
func (s sros) Dial(ctx context.Context) error {
	return establishConnection(ctx, s, s.Connection, s.basePrompt(), []string{"environment more false"})
}

// Enable elevates privileges.
func (s sros) Enable(ctx context.Context) error {
	return nil
}

// Run executes a command on a device.
func (s sros) Run(ctx context.Context, command string) (string, error) {
	output, err := s.RunUntil(ctx, command, s.basePrompt())
	if err != nil {
		return "", err
	}

	output = strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(output, "\n")
	result := ""
	// len-2 to cut off the context piece of the prompt aka [/]
	for i := 1; i < len(lines)-2; i++ {
		// skip empty lines that sros adds for visual separation of diff commands
		// as we add it manually
		if (i == len(lines)-3) && (lines[i] == "") {
			continue
		} else {
			result += lines[i] + "\n"
		}
	}

	return result, nil
}

// RunUntil executes a command and reads until the provided prompt.
func (s sros) RunUntil(ctx context.Context, command string, prompt *regexp.Regexp) (string, error) {
	err := s.Connection.Send(ctx, command)
	if err != nil {
		return "", fmt.Errorf("unable to send command to device: %w", err)
	}

	reader := s.Connection.Recv(ctx)
	output, err := readUntilPrompt(ctx, reader, prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (s sros) basePrompt() *regexp.Regexp {
	return regexp.MustCompile(`^[ABCD]:\S+@\S+#`)
}
