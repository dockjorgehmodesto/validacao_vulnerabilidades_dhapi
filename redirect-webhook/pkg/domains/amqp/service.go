package amqp

import (
	"strconv"

	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting/connectors"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/message_processors"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/rabbitmq"
	"conductor.com.br/pier/hooks/pkg/domains/event_controller"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Service interface {
	// CRUD Related
	Save(key *tenants.TenantKey, connection connecting.Connection, autoConnect bool) connecting.Connection
	ValidateIfEventTypeIsAvailable(key *tenants.TenantKey, eventType string)
	PartialUpdate(key *tenants.TenantKey, id string, values map[string]interface{}) connecting.Connection
	ListEnabledByIssuer(key *tenants.TenantKey) ([]connecting.Connection, error)
	FindById(key *tenants.TenantKey, id string) connecting.Connection
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) helpers.Paginator
	Delete(key *tenants.TenantKey, id string)
	DeleteByIssuer(key *tenants.TenantKey)

	// Connection Management and listening
	BroadcastConnect(tenantKey *tenants.TenantKey, connection connecting.Connection)
	BroadcastConnectById(tenantKey *tenants.TenantKey, id string)
	BroadcastDisconnect(tenantKey *tenants.TenantKey, connection connecting.Connection)
	BroadcastDisconnectByKey(tenantKey *tenants.TenantKey, key string)
	Get(key string) *connecting.Connection
	FlushAllForTenant(key *tenants.TenantKey)
	ListOpenConnections(filterByIssuerId int, filterByEventType string) []connectors.Connector
	broadcastRecheckConnections(tenantKey *tenants.TenantKey)
	broadcastFlushConnections(tenantKey *tenants.TenantKey)
	RecheckConnections(tenantKey *tenants.TenantKey)
	ReloadAndFlush(tenantKey *tenants.TenantKey)
	Disconnect(connectionId string) error
}

type ServiceImpl struct {
	repository              Repository
	connectionMap           ConnectionManager
	messageProcessorFactory message_processors.ProcessorFactory
	eventController         event_controller.Controller
}

func NewService(repository Repository, manager ConnectionManager, messageProcessorFactory message_processors.MessageProcessorFactory, controller event_controller.Controller) Service {
	return ServiceImpl{repository: repository, connectionMap: manager, messageProcessorFactory: messageProcessorFactory, eventController: controller}
}

func (s ServiceImpl) Save(key *tenants.TenantKey, connection connecting.Connection, autoConnect bool) connecting.Connection {
	s.ValidateIfEventTypeIsAvailable(key, connection.EventType)

	err := connection.ValidatePropertiesForBroker()
	Error{}.HandleInvalidBrokerProperties(err)

	processor, err := s.messageProcessorFactory.CreateMessageProcessor(connection, key)
	if err != nil {
		Error{}.HandleInvalidMessageProcessor(err)
	}
	connector, _ := connectors.CreateConnector(connection, key, processor)
	e := connector.Test()
	Error{}.HandleConnectionTestError(e)

	savedConnection, e := s.repository.Save(key, connection)
	Error{}.HandleDatabaseError(e)

	if autoConnect {
		s.BroadcastConnect(key, savedConnection)
	}

	_ = s.eventController.UpdateEventByIssuerAndName(key, key.IssuerId, connection.EventType)
	return savedConnection
}

// Connections' event type cannot have duplicate values
func (s ServiceImpl) ValidateIfEventTypeIsAvailable(key *tenants.TenantKey, eventType string) {
	_, e := s.repository.FindByEventType(key, eventType)
	Error{}.HandleEventTypeValidationOnCreate(e)
}

func (s ServiceImpl) Update(key *tenants.TenantKey, connection connecting.Connection) connecting.Connection {
	update, e := s.repository.Update(key, connection)
	Error{}.HandleDatabaseError(e)
	return update
}

func (s ServiceImpl) PartialUpdate(key *tenants.TenantKey, id string, values map[string]interface{}) connecting.Connection {
	connectionSearch := s.FindById(key, id)
	updatedConnection := s.ParseUpdateData(connectionSearch, values)

	err := updatedConnection.ValidatePropertiesForBroker()
	Error{}.HandleInvalidBrokerProperties(err)

	// This error will never occur as it is treated in the previous statement
	connection, _ := connectors.CreateConnector(updatedConnection, key, nil)
	testErr := connection.Test()
	Error{}.HandleConnectionTestError(testErr)

	_, e := s.repository.PartialUpdate(key, id, values)
	Error{}.HandleDatabaseErrorWithMessage(e, "error updating amqp connection")
	go s.broadcastFlushConnections(key)
	_ = s.eventController.UpdateEventByIssuerAndName(key, key.IssuerId, updatedConnection.EventType)
	return s.FindById(key, id)
}

func (s ServiceImpl) Delete(key *tenants.TenantKey, id string) {
	connection := s.FindById(key, id)
	e := s.repository.Delete(key, id)
	Error{}.HandleDatabaseErrorWithMessage(e, "error deleting amqp connection")
	_ = s.eventController.UpdateEventByIssuerAndName(key, key.IssuerId, connection.EventType)
	s.broadcastFlushConnections(key)
}

func (s ServiceImpl) ListEnabledByIssuer(key *tenants.TenantKey) ([]connecting.Connection, error) {
	return s.repository.ListEnabledByIssuer(key)
}

func (s ServiceImpl) FindById(key *tenants.TenantKey, id string) connecting.Connection {
	connection, e := s.repository.FindById(key, id)
	Error{}.HandleSearch(e)
	return connection
}

func (s ServiceImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) helpers.Paginator {
	paginator, e := s.repository.FindAllWithPagination(key, page, limit, sort, example)
	Error{}.HandlePagination(paginator, e)
	return paginator
}

func (s ServiceImpl) ConnectIfNotConnected(tenantKey *tenants.TenantKey, connection connecting.Connection) {
	if !s.connectionMap.IsConnected(connection.ID) {
		if canConnectBasedOnEnvironment(connection) {
			processor, err := s.messageProcessorFactory.CreateMessageProcessor(connection, tenantKey)
			if err != nil {
				Error{}.HandleInvalidMessageProcessor(err)
			}
			connector, _ := connectors.CreateConnector(connection, tenantKey, processor)
			connector.SetOnClose(s.connectionMap.DeleteByKey)
			s.connectionMap.ConnectIfNotConnected(connector)
		}
	}
}

func (s ServiceImpl) DeleteByIssuer(key *tenants.TenantKey) {
	e := s.repository.DeleteByIssuer(key)
	Error{}.HandleDatabaseErrorWithMessage(e, "error deleting amqp connection")
	s.eventController.UpdateEvents(key)
	s.broadcastFlushConnections(key)
}

// canConnectBasedOnEnvironment checks whether an connection can be established based on
// environment variables which will contain a list of issuers and/or a list of events.
// The feature can be toggled by using the environment variable "eventSelection.enabled"
func canConnectBasedOnEnvironment(connection connecting.Connection) bool {
	if viper.GetBool("eventSelection.enabled") {
		issuerEnabled := false
		for _, issuerIdString := range viper.GetStringSlice("eventSelection.issuers") {
			if issuerIdInt, err := strconv.Atoi(issuerIdString); err == nil && issuerIdInt == connection.IssuerId {
				issuerEnabled = true
			}
		}

		eventEnabled := false
		for _, event := range viper.GetStringSlice("eventSelection.events") {
			if event == connection.EventType {
				eventEnabled = true
			}
		}
		return issuerEnabled && eventEnabled
	}
	return true
}

func (s ServiceImpl) BroadcastConnect(tenantKey *tenants.TenantKey, connection connecting.Connection) {
	conn, _ := connectors.CreateConnector(connection, tenantKey, nil)
	e := conn.Test()
	Error{}.HandleConnectionTestError(e)
	if !connection.Enabled {
		connection.Enabled = true
		s.Update(tenantKey, connection)
	}
	s.broadcastRecheckConnections(tenantKey)
}

func (s ServiceImpl) BroadcastConnectById(tenantKey *tenants.TenantKey, id string) {
	connection := s.FindById(tenantKey, id)
	s.BroadcastConnect(tenantKey, connection)
}

func (s ServiceImpl) BroadcastDisconnect(tenantKey *tenants.TenantKey, connection connecting.Connection) {
	if connection.Enabled {
		connection.Enabled = false
		e := s.repository.UpdateStatus(tenantKey, connection)
		Error{}.HandleDatabaseErrorWithMessage(e, "error updating connection status")
	}
	s.broadcastRecheckConnections(tenantKey)
}

func (s ServiceImpl) BroadcastDisconnectByKey(tenantKey *tenants.TenantKey, key string) {
	connection := s.FindById(tenantKey, key)
	s.BroadcastDisconnect(tenantKey, connection)
}

func (s ServiceImpl) Get(key string) *connecting.Connection {
	return s.connectionMap.Get(key).GetConnection()
}

func (s ServiceImpl) FlushAllForTenant(key *tenants.TenantKey) {
	s.broadcastFlushConnections(key)
}

func (s ServiceImpl) ListOpenConnections(issuerId int, eventType string) []connectors.Connector {
	return s.connectionMap.ListOpenConnections(issuerId, eventType)
}

func (s ServiceImpl) broadcastRecheckConnections(tenantKey *tenants.TenantKey) {
	if !viper.GetBool("features.disableMessageReading") {
		rabbitmq.PublishString(tenantKey.ToString(), rabbitmq.ExchangeRecheckConnections, rabbitmq.QueueRecheckConnectionsKey)
	}
}

func (s ServiceImpl) broadcastFlushConnections(tenantKey *tenants.TenantKey) {
	if !viper.GetBool("features.disableMessageReading") {
		rabbitmq.PublishString(tenantKey.ToString(), rabbitmq.ExchangeConnectionFlush, rabbitmq.QueueConnectionsFlush)
	}
}

// Connection management
// Verifies for each enabled connection on database whether it should be connected or not, closing or opening connections if necessary.
func (s ServiceImpl) RecheckConnections(key *tenants.TenantKey) {
	if !viper.GetBool("features.disableMessageReading") {

		splunk.ConnectionManagementLog{}.SendInfo("processing recheck")

		enabledConnections, e := s.repository.ListEnabledByIssuer(key)
		if e != nil {
			Error{}.HandleListEnabledConnectionsOnRecheck(e)
		}
		s.addOnlyNewConnections(key, enabledConnections)
		s.disconnectDisabledConnectionsByIssuer(enabledConnections, key.IssuerId)
	}
}

func (s ServiceImpl) addOnlyNewConnections(key *tenants.TenantKey, enabledConnections []connecting.Connection) {
	for _, enabledConnection := range enabledConnections {
		s.ConnectIfNotConnected(key, enabledConnection)
	}
}

func (s ServiceImpl) disconnectDisabledConnectionsByIssuer(enabledConnections []connecting.Connection, issuerId int) {
	for _, openConnection := range s.connectionMap.ListOpenConnections(issuerId, "") {
		isEnabled := false
		for _, enabledConnection := range enabledConnections {
			if openConnection.GetId() == enabledConnection.ID {
				isEnabled = true
			}
		}

		if !isEnabled {
			e := s.Disconnect(openConnection.GetId())
			if e != nil {
				Error{}.HandleDisconnectOnRecheck(e)
			}
		}
	}
}

func (s ServiceImpl) ReloadAndFlush(tenantKey *tenants.TenantKey) {
	if !viper.GetBool("features.disableMessageReading") {

		splunk.ConnectionManagementLog{}.SendInfo("processing reset")

		s.cleanupConnectionsForIssuer(tenantKey.IssuerId)

		enabledConnections, e := s.repository.ListEnabledByIssuer(tenantKey)
		if e != nil {
			Error{}.HandleListEnabledConnectionsOnReset(e)
			return
		}

		for _, connection := range enabledConnections {
			if s.connectionMap.IsConnected(connection.ID) {
				err := s.Disconnect(connection.ID)
				if err != nil {
					splunk.ConnectionManagementLog{}.SendError("could not disconnect when flushing", err, &splunk.MessageOpts{
						IssuerId:     tenantKey.IssuerId,
						ConnectionId: connection.ID,
						EventType:    connection.EventType,
					})
					continue
				}
				s.ConnectIfNotConnected(tenantKey, connection)
				splunk.ConnectionManagementLog{}.SendInfo("successful connection reset")
			} else {
				s.ConnectIfNotConnected(tenantKey, connection)
			}
		}
	}
}

func (s ServiceImpl) cleanupConnectionsForIssuer(issuerId int) {
	cleanupErrors := make([]error, 0)
	openConnections := s.connectionMap.ListOpenConnections(issuerId, "")
	for _, connection := range openConnections {
		e := s.Disconnect(connection.GetId())
		if e != nil {
			cleanupErrors = append(cleanupErrors, e)
		}
	}
	if len(cleanupErrors) > 0 {
		for _, err := range cleanupErrors {
			splunk.ConnectionManagementLog{}.SendError("could not cleanup connection map", err, &splunk.MessageOpts{IssuerId: issuerId})
		}
	}
}

func (s ServiceImpl) ParseUpdateData(connection connecting.Connection, values map[string]interface{}) connecting.Connection {
	err := mapstructure.Decode(values, &connection)
	Error{}.HandleMapUpdate(err)
	return connection
}

// Disconnect
// Fetches a connection in the connection map, closes the connection, and deletes from the map
func (s ServiceImpl) Disconnect(connectionId string) error {
	connection := s.connectionMap.Get(connectionId)
	if connection == nil {
		return errors.New("connection not found in tenant map")
	}

	err := connection.Disconnect()
	if err != nil {
		return err
	}

	s.connectionMap.DeleteByKey(connectionId)
	return nil
}
