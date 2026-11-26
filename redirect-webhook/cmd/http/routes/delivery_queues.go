package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/delivery_queues"
	"net/http"

	"github.com/go-chi/chi"
)

type DeliveryQueues struct {
	service delivery_queues.Service
}

func NewDeliveryQueuesRouter(service delivery_queues.Service) DeliveryQueues {
	return DeliveryQueues{service: service}
}

func (ie DeliveryQueues) Route(r chi.Router) {
	r.Use(middleware.CheckIssuerId)
	r.Get("/", ie.findAll)
	r.Post("/", ie.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", ie.findById)
		r.Put("/", ie.update)
		r.Delete("/", ie.delete)
	})
	r.Route("/list", func(r chi.Router) {
		r.Post("/", ie.createList)
	})
	r.Route("/reset/{id}", func(r chi.Router) {
		r.Delete("/", ie.deleteAll)
	})
}

func (ie DeliveryQueues) findAll(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	var requestQueryParams delivery_queues.DeliveryQueueRequest
	helpers.BindQuery(r, &requestQueryParams)
	pages := ie.service.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort, requestQueryParams.ToMap())
	helpers.JsonResponse(w, delivery_queues.ConvertPagination(pages), http.StatusOK)
}

func (ie DeliveryQueues) create(w http.ResponseWriter, r *http.Request) {
	var input delivery_queues.DeliveryQueuePersist
	helpers.BindJson(r, &input)
	deliveryQueue := ie.service.Create(helpers.GetTenant(r), input.ToEntity())

	helpers.JsonResponse(w, delivery_queues.ConvertEntityToResponse(deliveryQueue), http.StatusCreated)
}

func (ie DeliveryQueues) findById(w http.ResponseWriter, r *http.Request) {
	deliveryQueueId := chi.URLParam(r, helpers.PathParamDefault)
	deliveryQueue := ie.service.FindById(helpers.GetTenant(r), deliveryQueueId)
	helpers.JsonResponse(w, delivery_queues.ConvertEntityToResponse(deliveryQueue), http.StatusOK)
}

func (ie DeliveryQueues) update(w http.ResponseWriter, r *http.Request) {
	deliveryQueueId := chi.URLParam(r, helpers.PathParamDefault)
	var input delivery_queues.DeliveryQueueUpdate
	helpers.BindJson(r, &input)
	deliveryQueue := ie.service.Update(helpers.GetTenant(r), deliveryQueueId, input.ToMap())
	helpers.JsonResponse(w, delivery_queues.ConvertEntityToResponse(deliveryQueue), http.StatusOK)
}

func (ie DeliveryQueues) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	ie.service.Delete(helpers.GetTenant(r), id)
	helpers.Ok(w, "")
}

func (ie DeliveryQueues) createList(w http.ResponseWriter, r *http.Request) {
	var input delivery_queues.DeliveryQueueList
	helpers.BindJson(r, &input)
	deliveryQueueList := ie.service.CreateList(helpers.GetTenant(r), input)
	helpers.JsonResponse(w, delivery_queues.ConvertListToResponse(deliveryQueueList), http.StatusCreated)
}

func (ie DeliveryQueues) deleteAll(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	ie.service.DeleteAll(helpers.GetTenant(r), id)
	helpers.Ok(w, "")
}
