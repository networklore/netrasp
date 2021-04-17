package netrasp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var errRead = errors.New("reader error")

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

type readResult struct {
	n   int
	err error
}

func (c *contextReader) Read(p []byte) (n int, err error) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	rrCh := make(chan *readResult)

	go func() {
		select {
		case rrCh <- c.read(p):
		case <-ctx.Done():
		}
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case rr := <-rrCh:
		return rr.n, rr.err
	}
}

func (c *contextReader) read(p []byte) *readResult {
	n, err := c.r.Read(p)

	return &readResult{n, err}
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
	r = newContextReader(ctx, r)
	for {
		buffer := make([]byte, 10000)

		bytes, err := r.Read(buffer)
		if err != nil {
			return "", fmt.Errorf("error reading output from device %w: %v", errRead, err)
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
