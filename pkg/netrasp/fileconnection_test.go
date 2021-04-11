package netrasp

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// fileConnection is a fake connection that reads input from files.
type fileConnection struct {
	command   string
	directory string
	host      *host
}

// Dial opens an SSH connection.
func (f *fileConnection) Dial(ctx context.Context) error {
	return nil
}

// GetHost returns information about the connected host.
func (f *fileConnection) GetHost() *host {
	return f.host
}

// Close disconnects from the device.
func (f *fileConnection) Close(ctx context.Context) error {
	return nil
}

// Send is used to write commands to the device.
func (f *fileConnection) Send(ctx context.Context, command string) error {
	f.command = command

	return nil
}

// Recv is used to read data from the device.
func (f *fileConnection) Recv(ctx context.Context) io.Reader {
	command := strings.ReplaceAll(f.command, " ", "_")
	content, err := ioutil.ReadFile(f.directory + "/" + command + ".txt")
	if err != nil {
		log.Fatalf("unable load output from command: %v", err)
	}
	output := string(content)
	output = strings.TrimRight(output, "\n")

	return newContextReader(ctx, strings.NewReader(output))
}

func withFakeFileConnection(directory string) ConfigOpt {
	return newFuncConfigOpt(func(c *config) {
		c.Connection = &fileConnection{directory: directory, command: "initial", host: c.Host}
	})
}
