package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/issuer_events"
	"github.com/go-chi/chi"
	"net/http"
)

type IssuerEvents struct {
	service issuer_events.Service
}

func NewIssuerEventsRouter(service issuer_events.Service) IssuerEvents {
	return IssuerEvents{service: service}
}

func (ie IssuerEvents) Route(r chi.Router) {
	r.Use(middleware.CheckIssuerId)
	r.Route("/events", func(r chi.Router) {
		r.Get("/", ie.findAll)
		r.Post("/", ie.create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ie.findById)
			r.Put("/", ie.update)
			r.Delete("/", ie.delete)

		})
	})
}

// @Title listAllIssuerEvents
// @Tags Issuer, Internals
// @Security issuerIdAuthentication
// @Summary List issuer events
// @Description List all issuer events
// @Param   page 				query    int	false "Page of the current result set"
// @Param   limit 				query    int	false "Amount of data that should be returned"
// @Param   sort 				query    string	false "Attribute to order the list"
// @Param   sort 				query    int	false "Attribute to order the list"
// @Param   issuerId			query    string	false "Identifier of the issuer"
// @Param   name				query    string	false "Name of the event"
// @Success 200 {object} issuer_events.IssuerEventResponsePage
// @Success 204 "No Content"
// @Accept json
// @Router /issuers/events [get]
func (ie IssuerEvents) findAll(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	var requestQueryParams issuer_events.IssuerEventRequest
	helpers.BindQuery(r, &requestQueryParams)
	pages := ie.service.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort, requestQueryParams.ToMap())
	helpers.JsonResponse(w, issuer_events.ConvertPagination(pages), http.StatusOK)
}

// @Title createIssuerEvents
// @Tags Issuer, Internals
// @Security issuerIdAuthentication
// @Summary Associate event to issuer. This will make the event available for the issuer.
// @Description Associates an event type to an issuer. This will make the event available for the issuer.\
// @Description An issuer can't be associated to two events of the same connectionId or with the same name.
// @Param   content body  issuer_events.IssuerEventPersist  true  "Object for associating the issuer an event"
// @Success 200 {object} issuer_events.IssuerEventResponse
// @Success 400 "event name is already taken"
// @Accept json
// @Router /issuers/events [post]
func (ie IssuerEvents) create(w http.ResponseWriter, r *http.Request) {
	var input issuer_events.IssuerEventPersist
	helpers.BindJson(r, &input)
	issuerEvent := ie.service.Create(helpers.GetTenant(r), input.ToEntity())

	helpers.JsonResponse(w, issuer_events.ConvertEntityToResponse(issuerEvent), http.StatusCreated)
}

// @Title findIssuerEventById
// @Tags Issuer, Internals
// @Security issuerIdAuthentication
// @Summary Find issuer event by id
// @Description Finds an issuer-event association by its identifier
// @Param id path string true "The identifier for the issuer-event association"
// @Success 200 {array} issuer_events.IssuerEventResponse
// @Success 204 "No Content"
// @Accept json
// @Router /issuers/events/{id} [get]
func (ie IssuerEvents) findById(w http.ResponseWriter, r *http.Request) {
	issuerEventId := chi.URLParam(r, helpers.PathParamDefault)
	issuerEvent := ie.service.FindById(helpers.GetTenant(r), issuerEventId)
	helpers.JsonResponse(w, issuer_events.ConvertEntityToResponse(issuerEvent), http.StatusOK)
}

// @Title updateIssuerEvent
// @Tags Issuer, Internals
// @Security issuerIdAuthentication
// @Summary Update issuer event by id
// @Description Updates an issuer-event association by its identifier
// @Param id path string true "The identifier for the issuer-event association"
// @Param   content body  issuer_events.IssuerEventUpdate  true  "Object for updating an issuer-event association"
// @Success 200 {array} issuer_events.IssuerEventResponse
// @Success 204 "No Content"
// @Success 400 "event name is already taken"
// @Accept json
// @Router /issuers/events/{id} [put]
func (ie IssuerEvents) update(w http.ResponseWriter, r *http.Request) {
	issuerEventId := chi.URLParam(r, helpers.PathParamDefault)
	var input issuer_events.IssuerEventUpdate
	helpers.BindJson(r, &input)
	issuerEvent := ie.service.Update(helpers.GetTenant(r), issuerEventId, input.ToMap())
	helpers.JsonResponse(w, issuer_events.ConvertEntityToResponse(issuerEvent), http.StatusOK)
}

// @Title deleteIssuerEvent
// @Tags Issuer, Internals
// @Security issuerIdAuthentication
// @Summary Delete issuer event association
// @Description Deletes an issuer-event association by its identifier
// @Param id path string true "The identifier for the issuer-event association"
// @Success 200 "Deleted"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /issuers/events/{id} [delete]
func (ie IssuerEvents) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	ie.service.Delete(helpers.GetTenant(r), id)
	helpers.Ok(w, "")
}
