package netrasp

import (
	"context"
	"io"
)

// Connection defines an interface for Netrasp connections.
type Connection interface {
	Dial(context.Context) error
	Close(context.Context) error
	Send(context.Context, string) error
	Recv(context.Context) io.Reader
	GetHost() *Host
}

// Host defines host specific information.
type Host struct {
	Address  string
	Port     int
	Platform Platform
	password string
}
