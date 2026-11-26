package webhooks

import (
	"time"

	"conductor.com.br/pier/hooks/pkg/helpers"
)

type WebHookPersist struct {
	EventType string          `json:"eventType" valid:"Required"`
	URL       string          `json:"url" valid:"Required; MaxSize(2000)"`
	Headers   []HeaderPersist `json:"headers"`
}

type HeaderPersist struct {
	Key   string `json:"key"`
	Value string `json:"value" valid:"MaxSize(2000)"`
}

func (h HeaderPersist) ConvertToEntity() Header {
	return Header{
		Key:   h.Key,
		Value: h.Value,
	}
}

func (s WebHookPersist) ConvertToEntity() WebHook {
	headerSlice := make([]Header, 0)
	for _, v := range s.Headers {
		headerSlice = append(headerSlice, v.ConvertToEntity())
	}

	return WebHook{
		EventType: s.EventType,
		URL:       s.URL,
		Headers:   headerSlice,
	}
}

type WebHookUpdate struct {
	EventType string          `json:"eventType" valid:"Required" model:"event_type"`
	URL       string          `json:"url" valid:"Required; MaxSize(2000)"`
	Headers   []HeaderPersist `json:"headers"`
}

type WebHookResponse struct {
	ID        string           `json:"id"`
	EventType string           `json:"eventType"`
	URL       string           `json:"url"`
	Headers   []HeaderResponse `json:"headers,omitempty"`
}

type HeaderResponse struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func ConvertHeaderToResponse(header Header) HeaderResponse {
	return HeaderResponse{
		ID:        header.ID,
		Key:       header.Key,
		Value:     header.Value,
		CreatedAt: header.CreatedAt,
		UpdatedAt: header.UpdatedAt,
		DeletedAt: header.DeletedAt,
	}
}

func ConvertEntityToResponse(hook WebHook) WebHookResponse {
	headerSlice := make([]HeaderResponse, 0)

	for _, v := range hook.Headers {
		headerSlice = append(headerSlice, ConvertHeaderToResponse(v))
	}

	return WebHookResponse{
		ID:        hook.ID,
		EventType: hook.EventType,
		URL:       hook.URL,
		Headers:   headerSlice,
	}
}

type WebHookResponsePage struct {
	helpers.PageResponse
	Records []WebHookResponse `json:"records"`
}

func ConvertPagination(paginator helpers.Paginator) WebHookResponsePage {
	slice := make([]WebHookResponse, 0)
	records := paginator.Records.(*[]WebHook)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToResponse(rec))
	}

	return WebHookResponsePage{
		PageResponse: helpers.PageResponse{
			TotalRecord: paginator.TotalRecord,
			TotalPage:   paginator.TotalPage,
			Offset:      paginator.Offset,
			Limit:       paginator.Limit,
			Page:        paginator.Page,
			PrevPage:    paginator.PrevPage,
			NextPage:    paginator.NextPage,
		},
		Records: slice,
	}
}

type WebHookTrigger struct {
	URL           string          `json:"url" valid:"Required"`
	Headers       []HeaderTrigger `json:"headers"`
	Data          string          `json:"data"`
	ForceMockData bool            `json:"forceMockData"`
	EventType     string          `json:"eventType"`
	PublicKey     string          `json:"publicKey"` // If informed, this will be used to encrypt the payload
}

func (w WebHookTrigger) ConvertToEntity() WebHook {
	headers := make([]Header, 0)
	for _, header := range w.Headers {
		headers = append(headers, header.ConvertToEntity())
	}

	return WebHook{
		EventType: w.EventType,
		URL:       w.URL,
		Headers:   headers,
	}
}

type WebHookTriggerResponse struct {
	Status          int    `json:"status"`
	Data            string `json:"data"`
	EncryptedAesKey string `json:"encryptedAesKey"`
}

type HeaderTrigger struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h HeaderTrigger) ConvertToEntity() Header {
	return Header{
		Key:   h.Key,
		Value: h.Value,
	}
}
