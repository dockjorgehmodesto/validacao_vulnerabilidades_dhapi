package webhooks

import (
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestWebHook_WithHeaders(t *testing.T) {
	t.Run("adding an empty map", func(t *testing.T) {
		webhook := WebHook{}
		webhook = webhook.WithHeader("a", "b")

		newMap := make(map[string]string)
		webhook = webhook.WithHeaders(newMap)

		assert.NotEmpty(t, webhook.Headers)
		assert.Equal(t, len(webhook.Headers), 1)
		assert.Equal(t, webhook.Headers[0].Key, "a")
		assert.Equal(t, webhook.Headers[0].Value, "b")
	})

	t.Run("adding a map with one value", func(t *testing.T) {
		webhook := WebHook{}
		webhook = webhook.WithHeader("a", "b")

		newMap := make(map[string]string)
		newMap["z"] = "x"
		webhook = webhook.WithHeaders(newMap)

		assert.NotEmpty(t, webhook.Headers)
		assert.Equal(t, len(webhook.Headers), 2)
	})

	t.Run("adding a map with multiple values", func(t *testing.T) {

		webhook := WebHook{}
		webhook = webhook.WithHeader("a", "b")

		newMap := make(map[string]string)
		randInt := (rand.Int() % 100) + 10
		for i := 0; i < randInt; i++ {
			newMap[random.String(uint8(randInt))] = random.String(uint8(randInt))
		}

		webhook = webhook.WithHeaders(newMap)

		assert.NotEmpty(t, webhook.Headers)
		assert.Equal(t, len(webhook.Headers), randInt+1)
	})

}
