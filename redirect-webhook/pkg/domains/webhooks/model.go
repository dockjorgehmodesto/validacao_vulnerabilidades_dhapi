package webhooks

import (
	"context"
	"net/http"
	"time"

	"conductor.com.br/pier/hooks/pkg/helpers"
)

type WebHook struct {
	helpers.GormHooks
	ID        string
	EventType string
	URL       string
	Headers   []Header
	IssuerId  int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Header struct {
	helpers.GormHooks
	ID        string
	WebHookId string
	Key       string
	Value     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (h Header) TableName() string {
	return "web_hook_headers"
}

func (w WebHook) GetHeaderMap() map[string]string {
	resultMap := make(map[string]string)
	for _, header := range w.Headers {
		resultMap[header.Key] = header.Value
	}
	return resultMap
}

func (w WebHook) Call(ctx context.Context, data string) (*helpers.HttpResponse, error) {
	return helpers.HttpRequest(ctx, http.DefaultClient, w.URL, w.GetHeaderMap(), []byte(data), http.MethodPost)
}

func (w WebHook) WithHeader(key string, value string) WebHook {
	newHeader := Header{
		Key:   key,
		Value: value,
	}
	w.Headers = append(w.Headers, newHeader)
	return w
}

// Copies a webhook with a new map of headers
func (w WebHook) WithHeaders(headersToAdd map[string]string) WebHook {
	for key, value := range headersToAdd {
		newHeader := Header{
			Key:   key,
			Value: value,
		}
		w.Headers = append(w.Headers, newHeader)
	}
	return w
}
