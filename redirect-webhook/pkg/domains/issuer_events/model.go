package issuer_events

import "time"

type IssuerEvent struct {
	Id           string
	Name         string
	IssuerId     int
	ConnectionId *string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
