package netrasp

import (
	"context"
)

type Asa struct {
	Ios
}

func (a Asa) Dial(ctx context.Context) error {
	return establishConnection(ctx, a, a.Connection, a.basePrompt(), "no terminal pager")
}
