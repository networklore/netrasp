package netrasp

import (
	"context"
	"io"
)

// connection defines an interface for Netrasp connections.
type connection interface {
	Dial(context.Context) error
	Close(context.Context) error
	Send(context.Context, string) error
	Recv(context.Context) io.Reader
	GetHost() *host
}

// host defines host specific information.
type host struct {
	Address  string
	Port     int
	Platform Platform
	password string
}
