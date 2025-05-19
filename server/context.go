package server

import (
	"context"
	"time"
)

func timeout(ctx context.Context, t time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, t)
}
