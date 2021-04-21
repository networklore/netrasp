package netrasp

import (
	"context"
	"fmt"
	"regexp"
)

func establishConnection(ctx context.Context, p Platform, c connection, prompt *regexp.Regexp, preparationCommands []string) error {
	err := c.Dial(ctx)
	if err != nil {
		return fmt.Errorf("unable to open connection: %w", err)
	}

	reader := c.Recv(ctx)
	// Make sure that we find the initial prompt to clear the buffer before we continue
	_, err = readUntilPrompt(ctx, reader, prompt)
	if err != nil {
		return fmt.Errorf("unable to find the initial prompt: %w", err)
	}

	for _, command := range preparationCommands {
		output, err := p.Run(ctx, command)
		if err != nil {
			return fmt.Errorf("unable to prepare session using the command '%s': %w", output, err)
		}
	}

	return nil
}
