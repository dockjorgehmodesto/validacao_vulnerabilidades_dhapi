package routes

import (
	"net/http"
	"strconv"

	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/pkg/domains/amqp"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"github.com/go-chi/chi"
)

type AMQP struct {
	service        amqp.Service
	tenantsService tenants.Service
}

func NewAMQPRouter(service amqp.Service, tenantsService tenants.Service) AMQP {
	return AMQP{service: service, tenantsService: tenantsService}
}

func (s AMQP) Route(r chi.Router) {
	r.Use(middleware.CheckIssuerId)

	r.Get("/", s.findAllConfigurationsWithPagination)
	r.Post("/", s.createConfiguration)
	r.Route("/{id}", func(idRoutes chi.Router) {
		idRoutes.Use(middleware.ValidateUUID(helpers.PathParamDefault))
		idRoutes.Get("/", s.findConfigurationById)
		idRoutes.Put("/", s.updateConfiguration)
		idRoutes.Delete("/", s.delete)

		idRoutes.Post("/connect", s.connect)
		idRoutes.Post("/disconnect", s.disconnect)

	})

	r.Route("/connections", func(connectionsRoutes chi.Router) {
		connectionsRoutes.Get("/", s.listOpenConnections)
		connectionsRoutes.Post("/flush", s.flushAllConnections)
	})

	r.Route("/issuer", func(issuerRoutes chi.Router) {
		issuerRoutes.Delete("/", s.deleteByIssuer)

	})

}

// @Title listConnections
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary List amqp connection data
// @Description List all the AMQP connection properties based on filters and pagination
// @Description These properties are used to manage connections and handle it's listeners
// @Param   page 				query    int      false "Page of the current result set"
// @Param   limit 				query    int      false "Amount of data that should be returned"
// @Param   sort 				query    string   false "Attribute to order the list (e.g. 'connectionString', 'connectionString DESC', 'connectionString ASC', 'id DESC, connectionString')"
// @Param   eventType			query    string   false "The type of event associated to the AMQP connection"
// @Param   enabled				query    bool	  false "Filters by the enabled flag"
// @Param   issuerId			query    int	  false "Filters by issuerId"
// @Success 200 {object} amqp.AMQPResponsePage
// @Success 204 "No Content"
// @Accept json
// @Router /amqp [get]
func (s AMQP) findAllConfigurationsWithPagination(w http.ResponseWriter, r *http.Request) {
	page := helpers.GetPagination(r)
	var requestQueryParams amqp.AMQPRequest
	helpers.BindQuery(r, &requestQueryParams)

	pages := s.service.FindAllWithPagination(helpers.GetTenant(r), page.Page, page.Limit, page.Sort, requestQueryParams.ToMap())

	helpers.JsonResponse(w, amqp.ConvertPagination(pages), http.StatusOK)
}

// @Title createConnection
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary Create amqp connection
// @Description Creates a new AMQP broker connection.\
// @Description Connection topics are unique per issuer, and must match one of the event types.\
// @Description If the autoConnect property is true, the application will inform all nodes and connect to it. Refer to the /amqp/{id}/connect method for more details on this process.\
// @Description \
// @Description \
// @Description Available processorType's: `"CDC_MESSAGE_PROCESSOR", "ODIN_MESSAGE_PROCESSOR"`.\
// @Description \
// @Description Available brokers: `"SERVICEBUS", "RABBITMQ"`.
// @Param   content body  amqp.AMQPPersist  true  "Object for persisting the AMQP server connection data"
// @Success 200 {object} amqp.AMQPResponse
// @Failure 400 "Bad request"
// @Accept json
// @Router /amqp [post]
func (s AMQP) createConfiguration(w http.ResponseWriter, r *http.Request) {
	var input amqp.AMQPPersist
	helpers.BindJson(r, &input)
	connection := s.service.Save(helpers.GetTenant(r), input.ConvertToEntity(), input.AutoConnect)

	helpers.JsonResponse(w, amqp.ConvertEntityToResponse(connection), http.StatusCreated)
}

// @Title findAMQPConnectionById
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary Find AMQP connection by ID
// @Description Finds for a AMQP connection with given ID
// @Param 		id 		path	string 		true 	"The identifier for the amqp connection properties"
// @Success 200 {object} amqp.AMQPResponse
// @Success 204 "No Content"
// @Accept json
// @Router /amqp/{id} [get]
func (s AMQP) findConfigurationById(w http.ResponseWriter, r *http.Request) {
	connectionId := chi.URLParam(r, helpers.PathParamDefault)
	connection := s.service.FindById(helpers.GetTenant(r), connectionId)
	helpers.JsonResponse(w, amqp.ConvertEntityToResponse(connection), http.StatusOK)
}

// @Title updateAMQPConnectionById
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary Update a AMQP connection
// @Description Only properties referring to the connections can be updated.
// @Param   id      path    string                true "The identifier for the AMQP connection"
// @Param   content body    amqp.AMQPUpdate      true "Object for updating a AMQP connection"
// @Success 200 {object} amqp.AMQPResponse
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/{id} [put]
func (s AMQP) updateConfiguration(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	var input amqp.AMQPUpdate
	helpers.BindJson(r, &input)
	connection := s.service.PartialUpdate(helpers.GetTenant(r), id, input.ToMap())
	helpers.JsonResponse(w, amqp.ConvertEntityToResponse(connection), http.StatusOK)
}

// @Title deleteAMQPConnectionById
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary Deletes an AMQP connection
// @Description Deleting an connection is irreversible and informs all clusters to disconnect to it
// @Param   id      path    string                true "The identifier for the AMQP connection"
// @Success 200 "Deleted"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/{id} [delete]
func (s AMQP) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	s.service.Delete(helpers.GetTenant(r), id)
	helpers.Ok(w, "")
}

// @Title connectById
// @Tags ConnectionManagement, Internals
// @Security issuerIdAuthentication
// @Summary Connects and listens to a AMQP topic
// @Description Broadcast all nodes that a new connection/listener is enabled.\
// @Description Verifies if the AMQP listener is not listening, then starts the connection procedure:\
// @Description The first step is to verify if there's a webhook that it's event type matches the listener's topic, associating it to the connection. If there's no webhook with the same name as the topic property, the listener will consume the messages anyways, as in the issuer does not want to receive it.\
// @Description Once the connection is set up with all the necessary data, it connects and starts listening. The listener's default behavior is to fetch the message from the queue (or topic, based on which broker is being used), encrypt it's payload and send it through the webhook.\
// @Description When listening, for each message the consumer will fetch the key vault having the latest expiration date\
// @Description For details on decrypting a message on the webhook's side, please refer to the "/encryption/key-vaults/{id}/aes/encrypt" resource documentation
// @Param   id      path    string                true "The identifier for the AMQP connection"
// @Success 200 "Connected"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/{id}/connect [post]
func (s AMQP) connect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	s.service.BroadcastConnectById(helpers.GetTenant(r), id)
	helpers.Success(w)
}

// @Title disconnectById
// @Tags ConnectionManagement, Internals
// @Security issuerIdAuthentication
// @Summary Disconnects an AMQP connection
// @Description Disables an AMQP listener and requests all nodes for disconnecting
// @Param   id      path    string                true "The identifier for the amqp"
// @Success 200 "Disconnected"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/{id}/disconnect [post]
func (s AMQP) disconnect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, helpers.PathParamDefault)
	s.service.BroadcastDisconnectByKey(helpers.GetTenant(r), id)
	helpers.Success(w)
}

// @Title listAllConnections
// @Tags ConnectionManagement, Internals
// @Security issuerIdAuthentication
// @Summary List open connections
// @Description Lists all active listeners
// @Param   eventType			query    string   false "The type of event associated to the AMQP connection"
// @Param   issuerId			query    int	  false "The issuer identification code"
// @Success 200 {array} amqp.AMQPResponse
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/connections [get]
func (s AMQP) listOpenConnections(w http.ResponseWriter, r *http.Request) {
	eventType := r.URL.Query().Get("eventType")
	issuerIdStr := r.URL.Query().Get("issuerId")
	issuerId := 0
	atoi, err := strconv.Atoi(issuerIdStr)
	if err == nil {
		issuerId = atoi
	}

	connections := s.service.ListOpenConnections(issuerId, eventType)

	response := make([]amqp.AMQPResponse, 0)
	for _, iConnection := range connections {
		response = append(response, amqp.ConvertEntityToResponse(*iConnection.GetConnection()))
	}
	helpers.Ok(w, response)
}

// @Title flushAllConnections
// @Tags ConnectionManagement, Internals
// @Security issuerIdAuthentication
// @Summary Flush connections
// @Description Flushes all AMQP connections for all issuers
// @Success 200 "Flushed"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/connections/flush [post]
func (s AMQP) flushAllConnections(w http.ResponseWriter, _ *http.Request) {
	tenantsSlice, _ := s.tenantsService.ListAll()
	for _, tenant := range tenantsSlice {
		tenantKey := tenant.GetTenantKey()
		s.service.FlushAllForTenant(&tenantKey)
	}
	helpers.Success(w)
}

// @Title deleteAMQPConnectionByIssuerId
// @Tags AMQP, Internals
// @Security issuerIdAuthentication
// @Summary Deletes an AMQP connection
// @Description Deleting an connection is irreversible and informs all clusters to disconnect to it
// @Param   id      path    string                true "The identifier for the AMQP connection"
// @Success 200 "Deleted"
// @Failure 204 "Connection not found"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Accept json
// @Router /amqp/issuer [delete]
func (s AMQP) deleteByIssuer(w http.ResponseWriter, r *http.Request) {
	// id := chi.URLParam(r, helpers.PathParamDefault)
	s.service.DeleteByIssuer(helpers.GetTenant(r))
	helpers.Ok(w, "")
}
