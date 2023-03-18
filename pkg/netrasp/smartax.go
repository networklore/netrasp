package netrasp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var generalSmartaxPrompt = regexp.MustCompile(`^[a-zA-Z0-9_()-]+[>|#]`)
var enableSmartaxPrompt = regexp.MustCompile(`^[a-zA-Z0-9_-]+#`)

// Smartax Huawei OLTs MA5XXX.
type smartax struct {
	Connection connection
}

// Close connection to device.
func (i smartax) Close(ctx context.Context) error {
	i.Connection.Close(ctx)

	return nil
}

// Configure device.
func (i smartax) Configure(ctx context.Context, commands []string) (ConfigResult, error) {
	var result ConfigResult
	_, err := i.Run(ctx, "config")
	if err != nil {
		return result, fmt.Errorf("unable to enter config mode: %w", err)
	}
	for _, command := range commands {
		output, err := i.Run(ctx, command)
		configCommand := ConfigCommand{Command: command, Output: output}
		result.ConfigCommands = append(result.ConfigCommands, configCommand)
		if err != nil {
			return result, fmt.Errorf("unable to run command '%s': %w", command, err)
		}
	}
	_, err = i.Run(ctx, "quit")
	if err != nil {
		return result, fmt.Errorf("unable to exit from config mode: %w", err)
	}

	return result, nil
}

// Dial opens a connection to a device.
func (i smartax) Dial(ctx context.Context) error {
	commands := []string{"undo smart", "scroll"}

	return establishConnection(ctx, i, i.Connection, i.basePrompt(), commands)
}

// Enable elevates privileges.
func (i smartax) Enable(ctx context.Context) error {
	_, err := i.RunUntil(ctx, "enable", enableSmartaxPrompt)
	if err != nil {
		return err
	}
	return nil
}

// Run executes a command on a device.
func (i smartax) Run(ctx context.Context, command string) (string, error) {
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
func (i smartax) RunUntil(ctx context.Context, command string, prompt *regexp.Regexp) (string, error) {
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

func (i smartax) basePrompt() *regexp.Regexp {
	return generalSmartaxPrompt
}
