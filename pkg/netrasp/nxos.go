package netrasp

import (
	"context"
)

type Nxos struct {
	Ios
}

func (n Nxos) Enable(ctx context.Context) error {
	return nil
}
