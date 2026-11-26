package message_processors

import (
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestMessage_MakeHeaders(t *testing.T) {
	t.Run("assert headers match", func(t *testing.T) {
		id := random.String(15) + "."
		data := random.String(15) + "."
		source := random.String(15) + "."
		deliveryCount := rand.Int() + 1

		message := Message{
			Id:            id,
			Data:          data,
			Source:        source,
			DeliveryCount: deliveryCount,
		}

		headers := message.MakeHeaders()

		assert.NotEmpty(t, headers[HeaderRequestId])
		assert.NotEmpty(t, headers[HeaderSource])

		assert.Equal(t, id, headers[HeaderRequestId])
		assert.Equal(t, source, headers[HeaderSource])
	})
}
