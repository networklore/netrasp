package netrasp

import (
	"context"
)

// Asa is the Netrasp driver for Cisco ASA.
type Asa struct {
	Ios
}

// Dial opens a connection to a device.
func (a Asa) Dial(ctx context.Context) error {
	return establishConnection(ctx, a, a.Connection, a.basePrompt(), "no terminal pager")
}
