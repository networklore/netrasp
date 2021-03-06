package netrasp

import (
	"context"
)

// Nxos is the Netrasp driver for Cisco NXOS.
type nxos struct {
	ios
}

// Enable just silently works on Nxos.
func (n nxos) Enable(ctx context.Context) error {
	return nil
}
