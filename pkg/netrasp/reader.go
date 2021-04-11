package netrasp

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (c *contextReader) Read(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, c.ctx.Err()
	default:
		return c.r.Read(p)
	}
}

func newContextReader(ctx context.Context, r io.Reader) io.Reader {
	return &contextReader{
		ctx: ctx,
		r:   r,
	}
}

// readUntilPrompt reads until the specified prompt is found and returns the read data.
func readUntilPrompt(ctx context.Context, r io.Reader, prompt *regexp.Regexp) (string, error) {
	var output string
	reader := newContextReader(ctx, r)

	for {
		buffer := make([]byte, 10000)
		time.Sleep(time.Millisecond * 10)
		bytes, err := reader.Read(buffer)
		if err != nil {
			return "", fmt.Errorf("error reading output from device: %w", err)
		}
		latestOutput := string(buffer[:bytes])

		output += latestOutput

		workingOutput := output
		workingOutput = strings.ReplaceAll(workingOutput, "\r\n", "\n")
		workingOutput = strings.ReplaceAll(workingOutput, "\r", "\n")
		lines := strings.Split(workingOutput, "\n")
		matches := prompt.FindStringSubmatch(lines[len(lines)-1])
		if len(matches) != 0 {
			break
		}
	}

	return output, nil
}
