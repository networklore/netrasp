package netrasp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var errRead = errors.New("reader error")

var minReadSize = 10000
var maxReadSize = 10 * 1000 * 1000

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
	var output = new(bytes.Buffer)
	rc := newContextReader(ctx, r)
	readSize := minReadSize
	for {
		b := make([]byte, readSize)

		n, err := rc.Read(b)
		if err != nil {
			return "", fmt.Errorf("error reading output from device %w: %v", errRead, err)
		}
		if readSize == n && readSize < maxReadSize {
			readSize *= 2
		}
		output.Write(b[:n])
		tempSlice := bytes.ReplaceAll(output.Bytes(), []byte("\r\n"), []byte("\n"))
		tempSlice = bytes.ReplaceAll(tempSlice, []byte("\r"), []byte("\n"))
		if prompt.Match(tempSlice[bytes.LastIndex(tempSlice, []byte("\n"))+1:]) {
			break
		}
	}

	return output.String(), nil
}
