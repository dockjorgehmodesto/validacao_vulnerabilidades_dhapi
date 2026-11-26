package timer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAppendTrackerToContext(t *testing.T) {
	t.Run("with empty context", func(t *testing.T) {
		ctx := context.Background()
		ctxWithTracker := AppendTrackerToContext(ctx)
		ctxValue := ctxWithTracker.Value(ctxKey)
		assert.NotNil(t, ctxValue)
		assert.NotPanics(t, func() {
			tracker := ctxValue.(ContextTimeTracker)
			assert.Empty(t, tracker)
		})

	})
	t.Run("with context already set", func(t *testing.T) {
		ctx := context.Background()
		ctxWithTracker := AppendTrackerToContext(ctx)
		ctxValue := ctxWithTracker.Value(ctxKey)

		trackerMap := ctxValue.(ContextTimeTracker)
		trackerMap["test"] = 1234

		// Overwriting should not change the map
		ctxWithTracker = AppendTrackerToContext(ctx)
		assert.NotNil(t, ctxValue)
		assert.NotPanics(t, func() {
			tracker := ctxValue.(ContextTimeTracker)
			assert.NotEmpty(t, tracker)
		})
	})
}

func TestGetTracker(t *testing.T) {
	t.Run("with empty context", func(t *testing.T) {
		ctx := context.Background()
		ctxWithTracker := AppendTrackerToContext(ctx)
		ctxValue := GetTracker(ctxWithTracker)
		assert.NotNil(t, ctxValue)
	})
	t.Run("with context already set", func(t *testing.T) {
		ctx := context.Background()
		ctxWithTracker := AppendTrackerToContext(ctx)
		ctxValue := ctxWithTracker.Value(ctxKey)

		trackerMap := ctxValue.(ContextTimeTracker)
		trackerMap["test"] = 1234

		ctxWithTracker = AppendTrackerToContext(ctxWithTracker)
		assert.NotNil(t, GetTracker(ctxWithTracker))
		assert.NotEmpty(t, GetTracker(ctxWithTracker))
	})
}

func TestTrack(t *testing.T) {
	t.Run("on context without ContextTimeTracker", func(t *testing.T) {
		ctx := context.Background()
		startedAt := time.Now()
		time.Sleep(100 * time.Millisecond)
		Track(ctx, "testkey", startedAt)
		assert.Nil(t, GetTracker(ctx))
	})
	t.Run("on context with ContextTimeTracker", func(t *testing.T) {
		ctx := context.Background()
		ctx = AppendTrackerToContext(ctx)
		startedAt := time.Now()
		time.Sleep(100 * time.Millisecond)
		Track(ctx, "testkey", startedAt)
		assert.NotNil(t, GetTracker(ctx))
		assert.NotEmpty(t, GetTracker(ctx))
	})
}
