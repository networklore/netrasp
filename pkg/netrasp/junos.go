package netrasp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Junos is the Netrasp driver for JunOS devices.
type junos struct {
	Connection connection
}

// Close connection to device.s
func (j junos) Close(ctx context.Context) error {
	j.Connection.Close(ctx)

	return nil
}

// Configure device.
func (j junos) Configure(ctx context.Context, commands []string) (string, error) {
	var output string
	_, err := j.Run(ctx, "configure exclusive")
	if err != nil {
		return "", fmt.Errorf("unable to enter exclusive edit mode: %w", err)
	}
	for _, command := range commands {
		result, err := j.Run(ctx, command)
		if err != nil {
			return output, fmt.Errorf("unable to run command '%s': %w", command, err)
		}
		output += result
	}
	_, err = j.Run(ctx, "commit")
	if err != nil {
		return output, fmt.Errorf("unable to commit configuration: %w", err)
	}

	_, err = j.Run(ctx, "quit")
	if err != nil {
		return output, fmt.Errorf("unable to quit from configuration mode: %w", err)
	}

	return output, nil
}

// Dial opens a connection to a device.
func (j junos) Dial(ctx context.Context) error {
	return establishConnection(ctx, j, j.Connection, j.basePrompt(), "set cli screen-length 0")
}

// Enable elevates privileges.
func (j junos) Enable(ctx context.Context) error {
	return nil
}

// Run executes a command on a device.
func (j junos) Run(ctx context.Context, command string) (string, error) {
	output, err := j.RunUntil(ctx, command, j.basePrompt())
	if err != nil {
		return "", err
	}

	output = strings.ReplaceAll(output, "\r\n", "\n")
	lines := strings.Split(output, "\n")
	result := ""

	for i := 1; i < len(lines)-1; i++ {
		result += lines[i] + "\n"
	}

	return result, nil
}

// RunUntil executes a command and reads until the provided prompt.
func (j junos) RunUntil(ctx context.Context, command string, prompt *regexp.Regexp) (string, error) {
	err := j.Connection.Send(ctx, command)
	if err != nil {
		return "", fmt.Errorf("unable to send command to device: %w", err)
	}

	reader := j.Connection.Recv(ctx)
	output, err := readUntilPrompt(ctx, reader, prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (j junos) basePrompt() *regexp.Regexp {
	return regexp.MustCompile(`^([a-zA-Z0-9_-]+)([>|#]|(@[a-zA-Z0-9_-]+[#|>]))`)
}
