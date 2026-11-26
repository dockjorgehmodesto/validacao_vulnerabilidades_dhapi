package routes

import (
	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/webhooks"
	"gopkg.in/jeevatkm/go-model.v1"
	"net/http"

	"github.com/go-chi/chi"
)

type WebHooks struct {
	service webhooks.Service
}

func NewWebhooksRouter(service webhooks.Service) WebHooks {
	return WebHooks{service: service}
}

func (wh WebHooks) Route(r chi.Router) {
	r.Group(func(issuerIdGroup chi.Router) {
		issuerIdGroup.Use(middleware.CheckIssuerId)
		issuerIdGroup.Get("/", wh.findAllWithPagination)
		issuerIdGroup.Post("/", wh.create)
		issuerIdGroup.Route("/{id}", func(idRoutes chi.Router) {
			idRoutes.Use(middleware.ValidateUUID(helpers.PathParamDefault))
			idRoutes.Get("/", wh.findById)
			idRoutes.Put("/", wh.update)
			idRoutes.Delete("/", wh.delete)
		})
	})
	r.Post("/trigger", wh.trigger)
}

// @Title listWebhooks
// @Tags WebHook
// @Security issuerIdAuthentication
// @Summary List webhooks
// @Description Displays all WebHooks registered on the API, based on filters and pagination.\
// @Description WebHooks resources can be used to create ("/webhooks" POST), update ("/webhooks/{id}" PUT) and test a connection to an external service, as in the "/webhooks/{id}/trigger" resource.\
// @Description WebHooks are used when processing messages served by AMQP. The messages are encrypted and sent to the related external service.
// @Param   page 				query    int      false "Page of the current result set"
// @Param   limit 				query    int      false "Amount of data that should be returned"
// @Param   sort 				query    string   false "Attribute to order the list (e.g. 'eventType', 'eventType DESC', 'status ASC', 'id DESC, status')"
// @Success 200 {object} webhooks.WebHookResponsePage
// @Success 204 "No Content"
// @Accept json
// @Router /webhooks [get]
func (wh WebHooks) findAllWithPagination(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)

	pages := wh.service.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort)

	helpers.JsonResponse(w, webhooks.ConvertPagination(pages), http.StatusOK)
}

// @Title createWebhook
// @Tags WebHook
// @Security issuerIdAuthentication
// @Summary Add webhook
// @Description Adds a new external service to the API with it's URL and necessary authentication passed through HTTP Headers to be used as a WebHook\
// @Description To make a WebHook get triggered through message brokers, make sure it's "eventType" matches the consumer's topic or queue name\
// @Description All WebHooks calls are performed using the HTTP Verb "POST", as there's the need to pass secure data through the payload
// @Param   content body  webhooks.WebHookPersist  true  "Object for persisting the webhook data"
// @Success 200 {object} webhooks.WebHookResponse
// @Failure 400 "Bad request"
// @Accept json
// @Router /webhooks [post]
func (wh WebHooks) create(w http.ResponseWriter, r *http.Request) {
	var input webhooks.WebHookPersist
	helpers.BindJson(r, &input)
	webHook := wh.service.Create(helpers.GetTenant(r), input.ConvertToEntity())

	helpers.JsonResponse(w, webhooks.ConvertEntityToResponse(webHook), http.StatusCreated)
}

// @Title findWebhookById
// @Tags WebHook
// @Security issuerIdAuthentication
// @Summary Find webhook by ID
// @Description Finds for a webhook with given ID
// @Param 		id 		path	string 		true 	"The identifier for the webhook"
// @Success 200 {object} webhooks.WebHookResponse
// @Success 204 "No Content"
// @Accept json
// @Router /webhooks/{id} [get]
func (wh WebHooks) findById(w http.ResponseWriter, r *http.Request) {
	webHookId := chi.URLParam(r, helpers.PathParamDefault)
	webHook := wh.service.FindByIdWithPreload(helpers.GetTenant(r), webHookId)
	helpers.JsonResponse(w, webhooks.ConvertEntityToResponse(webHook), http.StatusOK)
}

// @Title updateWebHookById
// @Tags WebHook
// @Security issuerIdAuthentication
// @Summary Update a webhook
// @Description Updates a webhook
// @Param   id      path    string                true "The identifier for the webhook"
// @Param   content body    webhooks.WebHookUpdate      true "Object for updating a webhook"
// @Success 200 {object} webhooks.WebHookResponse
// @Failure 204 "WebHook not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /webhooks/{id} [put]
func (wh WebHooks) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	var input webhooks.WebHookUpdate
	helpers.BindJson(r, &input)
	updateMap, _ := model.Map(input)
	webHook := wh.service.Update(helpers.GetTenant(r), id, updateMap)
	helpers.JsonResponse(w, webhooks.ConvertEntityToResponse(webHook), http.StatusOK)
}

// @Title deleteWebHookById
// @Tags WebHook
// @Security issuerIdAuthentication
// @Summary Delete a webhook
// @Description Whenever you no longer want to receive messages on your API, simply delete it
// @Param   id      path    string                true "The identifier for the webhook"
// @Success 200 "Deleted"
// @Failure 204 "WebHook not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /webhooks/{id} [delete]
func (wh WebHooks) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	wh.service.Delete(helpers.GetTenant(r), id)
	helpers.JsonResponse(w, "", http.StatusOK)
}

// @Title triggerWebHookById
// @Tags WebHook
// @Summary Trigger a webhook
// @Description Fetches a webHook data, calling it's external service through an HTTP Post\
// @Description If the "forceMockData" property is set to true, the data will be overridden to a mock, according to the webhook's event type\
// @Description If you want to test the application's encryption cycle, try informing your *publicKey* in the request's body, and the payload will be encrypted. For decrypting it, make sure to use the key given either in the request body and the HTTP header. With the encrypted AES key, decrypt it with your private RSA key, then decrypt the payload received in your API.
// @Param   content body    webhooks.WebHookTrigger      true "Object for calling a webhook"
// @Success 200 {object} webhooks.WebHookTriggerResponse
// @Failure 204 "WebHook not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /webhooks/trigger [post]
func (wh WebHooks) trigger(w http.ResponseWriter, r *http.Request) {
	var input webhooks.WebHookTrigger
	helpers.BindJson(r, &input)
	status, responseBody, encryptedAesKey := wh.service.Trigger(r.Context(), input.ConvertToEntity(), input.Data, input.ForceMockData, input.EventType, input.PublicKey)
	helpers.JsonResponse(w, webhooks.WebHookTriggerResponse{
		Status:          status,
		Data:            responseBody,
		EncryptedAesKey: encryptedAesKey,
	}, http.StatusOK)
}
