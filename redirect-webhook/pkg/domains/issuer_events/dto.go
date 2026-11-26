package issuer_events

import (
	"conductor.com.br/pier/hooks/pkg/helpers"
)

type IssuerEventRequest struct {
	Name         *string `json:"name" model:"name,omitempty"`
	IssuerId     *int    `json:"issuerId,string" model:"issuer_id,omitempty"`
	ConnectionId *string `json:"connectionId" model:"connection_id,omitempty"`
}

func (r IssuerEventRequest) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(r)
}

type IssuerEventResponse struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	IssuerId     int     `json:"issuerId"`
	ConnectionId *string `json:"connectionId"`
}

func ConvertEntityToResponse(event IssuerEvent) IssuerEventResponse {
	return IssuerEventResponse{
		Id:           event.Id,
		Name:         event.Name,
		IssuerId:     event.IssuerId,
		ConnectionId: event.ConnectionId,
	}
}

type IssuerEventResponsePage struct {
	helpers.PageResponse
	Records []IssuerEventResponse `json:"records"`
}

func ConvertPagination(paginator helpers.Paginator) IssuerEventResponsePage {
	slice := make([]IssuerEventResponse, 0)
	records := paginator.Records.(*[]IssuerEvent)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToResponse(rec))
	}

	return IssuerEventResponsePage{
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

type IssuerEventNameResponse struct {
	Name string `json:"name"`
}

func ConvertEntityToNameResponse(event IssuerEvent) IssuerEventNameResponse {
	return IssuerEventNameResponse{
		Name: event.Name,
	}
}

type IssuerEventNameResponsePage struct {
	helpers.PageResponse
	Records []IssuerEventNameResponse `json:"records"`
}

func ConvertNamePagination(paginator helpers.Paginator) IssuerEventNameResponsePage {
	slice := make([]IssuerEventNameResponse, 0)
	records := paginator.Records.(*[]IssuerEvent)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToNameResponse(rec))
	}

	return IssuerEventNameResponsePage{
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

type IssuerEventPersist struct {
	Name     string `json:"name" valid:"Required;MaxSize(250)"`
	IssuerId int    `json:"issuerId" valid:"Required"`
}

func (i IssuerEventPersist) ToEntity() IssuerEvent {
	return IssuerEvent{
		Name:     i.Name,
		IssuerId: i.IssuerId,
	}
}

type IssuerEventUpdate struct {
	Name         string `json:"name" valid:"Required;MaxSize(250)" model:"name"`
	IssuerId     int    `json:"issuerId" valid:"Required" model:"issuer_id"`
	ConnectionId string `json:"connectionId" valid:"Required;Match(/[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}/)" model:"connection_id"`
}

func (i IssuerEventUpdate) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(i)
}
