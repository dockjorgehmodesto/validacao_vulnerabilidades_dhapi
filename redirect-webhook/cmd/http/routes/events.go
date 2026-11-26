package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/amqp"
	"conductor.com.br/pier/hooks/pkg/domains/issuer_events"
	"github.com/go-chi/chi"
	"net/http"
)

type Events struct {
	service             amqp.Service
	issuerEventsService issuer_events.Service
}

func NewEventTypesRouter(service amqp.Service, issuerEventsService issuer_events.Service) Events {
	return Events{service: service, issuerEventsService: issuerEventsService}
}

func (e Events) Route(r chi.Router) {
	r.Use(middleware.CheckIssuerId)
	r.Get("/types", e.findAll)
}

// @Title listEventTypes
// @Tags EventType
// @Security issuerIdAuthentication
// @Summary List event types
// @Description Displays all available events for the issuer, based on filters and pagination.
// @Param   page 				query    int      false "Page of the current result set"
// @Param   limit 				query    int      false "Amount of data that should be returned"
// @Param   sort 				query    string   false "Attribute to order the list (e.g. 'eventType', 'eventType DESC', 'status ASC', 'id DESC, status')"
// @Success 200 {object} issuer_events.IssuerEventNameResponsePage
// @Success 204 "No Content"
// @Accept json
// @Router /events/types [get]
func (e Events) findAll(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	pagination := e.issuerEventsService.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort, map[string]interface{}{"issuer_id": helpers.GetIssuerId(r)})

	helpers.JsonResponse(w, issuer_events.ConvertNamePagination(pagination), http.StatusOK)
}
