package message_processors

import (
	"strconv"
	"time"
)

const (
	HeaderRequestId     = "DataHubMessageId"
	HeaderGroupId       = "DataHubMessageGroupId"
	HeaderSource        = "DataHubMessageSource"
	HeaderDeliveryCount = "DataHubMessageDeliveryCount"
	HeaderSentAt        = "sentAt"
	HeaderGeneratedAt   = "generatedAt"
	HeaderPublishAt     = "publishAt"
	HeaderPrimaryKey    = "primaryKey"
	HeaderOperation     = "operation"
)

const (
	MessageSourceDefault         = "DEFAULT"
	MessageSourceDeadLetterQueue = "DEAD LETTER QUEUE"
)

type Message struct {
	Id            string
	GroupId       string
	Data          string
	Source        string
	DeliveryCount int
	EventType     string
	GeneratedAt   *time.Time
	PublishAt     *time.Time
	// The primary key of the event
	PrimaryKey uint64
	// Operation identifies if the event it's an insert, update or a delete
	Operation string
}

func (m Message) MakeHeaders() map[string]string {
	result := make(map[string]string)
	result[HeaderRequestId] = m.Id
	result[HeaderSource] = m.Source
	result[HeaderGroupId] = m.GroupId
	result[HeaderPrimaryKey] = strconv.Itoa(int(m.PrimaryKey))
	result[HeaderOperation] = m.Operation
	return result
}
