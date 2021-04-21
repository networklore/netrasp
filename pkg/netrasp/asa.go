package netrasp

import (
	"context"
)

// Asa is the Netrasp driver for Cisco ASA.
type asa struct {
	ios
}

// Dial opens a connection to a device.
func (a asa) Dial(ctx context.Context) error {
	return establishConnection(ctx, a, a.Connection, a.basePrompt(), []string{"no terminal pager"})
}
