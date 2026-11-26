package amqp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_HandleConnectionNotFoundByEventType(t *testing.T) {
	t.Run("throw error", func(t *testing.T) {
		assert.Panics(t, func() {
			Error{}.HandleActiveConnectionNotFoundByEventType(errors.New(""))
		})
	})

	t.Run("nil error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Error{}.HandleActiveConnectionNotFoundByEventType(nil)
		})
	})
}

func TestError_HandleResend(t *testing.T) {
	t.Run("throw error", func(t *testing.T) {
		assert.Panics(t, func() {
			Error{}.HandleResend(errors.New(""))
		})
	})

	t.Run("nil error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Error{}.HandleResend(nil)
		})
	})
}

func TestError_HandleEventWithoutEncryptionKey(t *testing.T) {
	t.Run("throw error", func(t *testing.T) {
		assert.Panics(t, func() {
			Error{}.HandleEventWithoutEncryptionKey(errors.New(""))
		})
	})

	t.Run("nil error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Error{}.HandleEventWithoutEncryptionKey(nil)
		})
	})
}

func TestError_HandleEventWithoutWebHook(t *testing.T) {
	t.Run("throw error", func(t *testing.T) {
		assert.Panics(t, func() {
			Error{}.HandleEventWithoutWebHook(errors.New(""))
		})
	})

	t.Run("nil error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Error{}.HandleEventWithoutWebHook(nil)
		})
	})
}
