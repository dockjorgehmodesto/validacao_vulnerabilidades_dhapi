package delivery

import "time"

type WebHookMessageDelivery struct {
	Payload   string            `json:"payload"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Metadata  Metadata          `json:"metadata"`
	IssuerId  int               `json:"issuerId"`
	Processor ProcessorTrace    `json:"processor"`
}

type Metadata struct {
	MessageId    string     `json:"messageId"`
	InternalId   string     `json:"internalId"`
	KeyVaultId   string     `json:"keyVaultId"`
	WebHookId    string     `json:"webHookId"`
	ConnectionId string     `json:"connectionId"`
	EventType    string     `json:"eventType"`
	BrokerType   string     `json:"brokerType"`
	ProcessedAt  time.Time  `json:"processedAt"`
	Host         string     `json:"host"`
	GroupId      string     `json:"groupId"`
	GeneratedAt  *time.Time `json:"generatedAt"`
	PrimaryKey   uint64     `json:"primaryKey"`
	Operation    string     `json:"operation"`
	PublishAt    *time.Time `json:"publishAt"`
}

func (m Metadata) GetOriginDelayMs() int64 {
	if m.GeneratedAt != nil {
		return time.Since(*m.GeneratedAt).Milliseconds()
	}
	return 0
}

func (m Metadata) GetPublishDelayMs() int64 {
	if m.PublishAt != nil {
		return time.Since(*m.PublishAt).Milliseconds()
	}
	return 0
}

type ProcessorTrace struct {
	InsertedOnDate time.Time              `json:"insertedOnDate"`
	DurationSec    int64                  `json:"durationSec"`
	Durations      map[string]interface{} `json:"durations"`
}
