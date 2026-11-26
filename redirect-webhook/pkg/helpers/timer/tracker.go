package timer

import (
	"context"
	"time"
)

const ctxKey = "timeTracker"

type ContextTimeTracker map[string]interface{}

func AppendTrackerToContext(ctx context.Context) context.Context {
	if value := ctx.Value(ctxKey); value == nil {
		return context.WithValue(ctx, ctxKey, make(ContextTimeTracker))
	}
	return ctx
}

func Track(ctx context.Context, key string, startTime time.Time) {
	ctxValue := ctx.Value(ctxKey)
	if ctxValue != nil {
		if tracker, ok := ctxValue.(ContextTimeTracker); ok {
			tracker[key] = time.Since(startTime).Milliseconds()
		}
	}
}

func GetTracker(ctx context.Context) ContextTimeTracker {
	ctxValue := ctx.Value(ctxKey)
	if ctxValue != nil {
		if tracker, ok := ctxValue.(ContextTimeTracker); ok {
			return tracker
		}
	}
	return nil
}
