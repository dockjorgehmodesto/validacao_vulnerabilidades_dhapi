package rabbitmq

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

const (
	defaultHeartbeat = 10 * time.Second
	defaultLocale    = "en_US"
)

var manager RabbitManager

type RabbitManager struct {
	connection     *Connection
	channel        *Channel
	isReconnecting bool
	reconnectMutex sync.Mutex
	actions        ListenerActions
}

func InitConnections(actions ListenerActions) {

	manager = RabbitManager{
		isReconnecting: false,
		actions:        actions,
		reconnectMutex: sync.Mutex{},
	}
	manager.start()
}

func (r *RabbitManager) start() {
	r.createConnection()
	r.connectToChannel()
	r.declareExchangesToChannel()
	r.declareQueuesToChannel()
	go r.initReceivers()
}

func (r *RabbitManager) restart() {
	if !r.isReconnecting {
		r.setReconnect(true)
		defer r.setReconnect(false)
		reconnectDelay := viper.GetDuration("rabbit.sleep")
		time.Sleep(reconnectDelay * time.Second)
		r.start()
	}
}

func (r *RabbitManager) setReconnect(b bool) {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()
	r.isReconnecting = b
}

func (r *RabbitManager) createConnection() {
	if r.connection != nil && !r.connection.IsClosed() {
		e := r.connection.Close()
		if e != nil {
			logging.Logger.Error("[RABBIT] could not close an open connection", zap.Error(e))
		}
	}

	var e error
	url := global.Viper.GetString("rabbit.dial")
	vHost := global.Viper.GetString("rabbit.virtualHost")
	r.connection, e = Dial(url, amqp.Config{Heartbeat: defaultHeartbeat, Locale: defaultLocale, Vhost: vHost})
	if e != nil {
		logging.Logger.Error("[RABBIT] could not connect to dial", zap.Error(e))
		r.restart()
	}
}

func (r *RabbitManager) connectToChannel() {
	if r.channel != nil && !r.channel.IsClosed() {
		e := r.channel.Close()
		if e != nil {
			logging.Logger.Error("[RABBIT] could not close an open channel", zap.Error(e))
		}
	}

	var e error
	r.channel, e = r.connection.Channel()
	if e != nil {
		logging.Logger.Error("[RABBIT] could not create channel", zap.Error(e))
		r.restart()
	}

	e = r.channel.Qos(viper.GetInt("rabbit.qos.prefetchCount"), viper.GetInt("rabbit.qos.prefetchSize"), true)
	if e != nil {
		logging.Logger.Error("[RABBIT] could not apply qos", zap.Error(e))
		r.restart()
	}
}

func formatQueueName(name string) string {
	s, _ := os.Hostname()
	port := global.Viper.GetString("server.port")
	return fmt.Sprintf("%s_%s_%s", name, s, port)
}

func (r *RabbitManager) declareQueuesToChannel() {
	for _, queue := range Queues() {
		name := queue.Name
		if queue.IsBroadcast {
			name = formatQueueName(queue.Name)
		}

		q, err := r.channel.QueueDeclare(
			name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args,
		)
		failOnError(err, "Failed to declare a queue: "+queue.Name)

		err = r.channel.QueueBind(
			name,
			queue.Key,
			queue.Exchange,
			queue.NoWait,
			queue.Args)

		failOnError(err, "Failed to bind a queue: "+q.Name)
	}
}

func (r *RabbitManager) declareExchangesToChannel() {
	for _, exchange := range Exchanges() {

		err := r.channel.ExchangeDeclare(
			exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.Args,
		)
		failOnError(err, "Failed to declare an exchange: "+exchange.Name)
	}
}

func (r *RabbitManager) initReceivers() {
	for _, listener := range r.actions.Build() {
		listener.ConsumeMessages()
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		logging.Logger.Error(fmt.Sprintf("[RABBIT] %s", msg), zap.Error(err))
	}
}
