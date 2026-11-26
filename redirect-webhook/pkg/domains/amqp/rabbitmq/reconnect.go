package rabbitmq

import (
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"time"

	"sync/atomic"

	"github.com/streadway/amqp"
)

// Debug debug log flag
var Debug bool = true

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Channel() (*Channel, error) {
	if c == nil || c.Connection == nil {
		return nil, errors.New("nil connection when creating channel")
	}
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			if ok {
				log.Error(reason)
				go manager.restart()
				return
			}
		}
	}()

	return channel, nil
}

// Dial wrap amqp.Dial, dial and get a reconnect connection
func Dial(url string, config amqp.Config) (*Connection, error) {
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			if ok {
				log.Error(reason)
				go manager.restart()
				return
			}
		}
	}()

	return connection, nil
}

// Channel amqp.Channel wapper
type Channel struct {
	*amqp.Channel
	closed int32
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return (atomic.LoadInt32(&ch.closed) == 1)
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume warp amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			if !ch.IsClosed() {

				d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
				if err != nil {
					if amqpErr, ok := err.(*amqp.Error); ok {
						log.Error("error when consuming messages", amqpErr)
						go manager.restart()
						return
					}
					debugf("consume failed, err: %v", err)
					continue
				}

				for msg := range d {
					deliveries <- msg
				}
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			reconnectDelay := viper.GetDuration("rabbit.sleep")
			time.Sleep(reconnectDelay * time.Second)
		}
	}()

	return deliveries, nil
}

func debugf(format string, args ...interface{}) {
	if !Debug {
		return
	}
	log.Printf(format, args...)
}
