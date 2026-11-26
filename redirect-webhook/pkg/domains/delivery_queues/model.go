package delivery_queues

import "time"

var EventTypeList = []string{
	"account", "address", "adjustment", "bankslip", "batch", "card", "creditanalysis",
	"declined", "limit", "payment", "person", "phone", "purchase", "purchase_approved",
	"purchase_canceled", "purchase_cleared", "purchase_processed", "standin", "statement",
	"token", "tokenized_purchase", "transaction"}

type DeliveryQueue struct {
	Id           string
	IssuerId     int
	EventType    string
	ExchangeName string
	QueueName    string
	QueueKey     string
	Enabled      bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
