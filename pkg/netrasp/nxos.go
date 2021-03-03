package netrasp

import (
	"context"
)

// Nxos is the Netrasp driver for Cisco NXOS.
type Nxos struct {
	Ios
}

// Enable just silently works on Nxos.
func (n Nxos) Enable(ctx context.Context) error {
	return nil
}
