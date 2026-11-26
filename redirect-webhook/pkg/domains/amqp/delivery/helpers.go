package delivery

import "github.com/streadway/amqp"

func getDeliveryCount(message amqp.Delivery) int {
	// If the retry header exists, it means that a message was already delivered once.
	// For fetching the delivery count, we get it's value and add 2, as the library starts counting from 0,
	// and before creating the header, the first deliver attempt was already made
	if mapVal, ok := message.Headers["retries"]; ok {
		return int(mapVal.(int32)) + 1
	}

	// If there's no retry header, it's the first delivery attempt
	return 1
}
