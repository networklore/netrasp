package netrasp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var generalPrompt = regexp.MustCompile(`(^[a-zA-Z0-9_-]+[#|>])|(^[a-zA-Z0-9_-]+\([a-z-]+\)#)`)
var enablePrompt = regexp.MustCompile(`^[Pp]assword:`)

// Ios is the Netrasp driver for Cisco IOS.
type ios struct {
	Connection connection
}

// Close connection to device.
func (i ios) Close(ctx context.Context) error {
	i.Connection.Close(ctx)

	return nil
}

// Configure device.
func (i ios) Configure(ctx context.Context, commands []string) (string, error) {
	var output string
	_, err := i.Run(ctx, "configure terminal")
	if err != nil {
		return "", fmt.Errorf("unable to enter config mode: %w", err)
	}
	for _, command := range commands {
		result, err := i.Run(ctx, command)
		if err != nil {
			return output, fmt.Errorf("unable to run command '%s': %w", command, err)
		}
		output += result
	}
	_, err = i.Run(ctx, "end")
	if err != nil {
		return output, fmt.Errorf("unable to exit from config mode: %w", err)
	}

	return output, nil
}

// Dial opens a connection to a device.
func (i ios) Dial(ctx context.Context) error {
	return establishConnection(ctx, i, i.Connection, i.basePrompt(), "terminal length 0")
}

// Enable elevates privileges.
func (i ios) Enable(ctx context.Context) error {
	_, err := i.RunUntil(ctx, "enable", enablePrompt)
	if err != nil {
		return err
	}
	host := i.Connection.GetHost()
	_, err = i.Run(ctx, host.password)

	if err != nil {
		return err
	}

	return nil
}

// Run executes a command on a device.
func (i ios) Run(ctx context.Context, command string) (string, error) {
	output, err := i.RunUntil(ctx, command, i.basePrompt())
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
func (i ios) RunUntil(ctx context.Context, command string, prompt *regexp.Regexp) (string, error) {
	err := i.Connection.Send(ctx, command)
	if err != nil {
		return "", fmt.Errorf("unable to send command to device: %w", err)
	}

	reader := i.Connection.Recv(ctx)
	output, err := readUntilPrompt(ctx, reader, prompt)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (i ios) basePrompt() *regexp.Regexp {
	return generalPrompt
}
