package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting/connectors"
	"github.com/pkg/errors"
	"sync"
)

type ConnectionManager interface {
	Connect(connection connectors.Connector)
	ConnectIfNotConnected(connection connectors.Connector)
	DeleteByKey(key string)
	Get(key string) connectors.Connector
	ListOpenConnections(issuerId int, eventType string) []connectors.Connector
	IsConnected(key string) bool
	FindByIssuerAndEventType(id int, eventType string) (connectors.Connector, error)
}

type ConnectionManagerImpl struct {
	openConnections map[string]connectors.Connector
	mutex           *sync.RWMutex
}

func NewConnectionManager() ConnectionManager {
	return &ConnectionManagerImpl{
		openConnections: make(map[string]connectors.Connector),
		mutex:           &sync.RWMutex{},
	}
}

func (c *ConnectionManagerImpl) Connect(connection connectors.Connector) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	connection.SetOnClose(c.DeleteByKey)
	c.openConnections[connection.GetId()] = connection
	go connection.Listen()
}

func (c *ConnectionManagerImpl) ConnectIfNotConnected(connection connectors.Connector) {
	if !c.IsConnected(connection.GetId()) {
		c.Connect(connection)
	}
}

func (c *ConnectionManagerImpl) IsConnected(key string) bool {
	return c.Get(key) != nil
}

func (c *ConnectionManagerImpl) DeleteByKey(key string) {
	connection := c.Get(key)
	if connection != nil {
		c.delete(connection)
	}
}

func (c *ConnectionManagerImpl) delete(connection connectors.Connector) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.openConnections, connection.GetId())
	connection = nil
	return
}

func (c *ConnectionManagerImpl) Get(key string) connectors.Connector {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if connection, ok := c.openConnections[key]; !ok {
		return nil
	} else {
		return connection
	}
}

func (c *ConnectionManagerImpl) ListOpenConnections(issuerId int, eventType string) []connectors.Connector {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	result := make([]connectors.Connector, 0)
	for _, v := range c.openConnections {
		if v != nil {
			result = append(result, v)
		}
	}

	if issuerId != 0 {
		result = c.filterConnectionsByIssuerId(result, issuerId)
	}

	if eventType != "" {
		result = c.filterConnectionsByEventType(result, eventType)
	}

	return result
}

func (c *ConnectionManagerImpl) filterConnectionsByIssuerId(input []connectors.Connector, issuerId int) []connectors.Connector {
	result := make([]connectors.Connector, 0)

	for _, connection := range input {
		if connection.GetIssuerId() == issuerId {
			result = append(result, connection)
		}
	}
	return result
}

func (c *ConnectionManagerImpl) filterConnectionsByEventType(input []connectors.Connector, eventType string) []connectors.Connector {
	result := make([]connectors.Connector, 0)

	for _, connection := range input {
		if connection.GetEventType() == eventType {
			result = append(result, connection)
		}
	}
	return result
}

func (c *ConnectionManagerImpl) FindByIssuerAndEventType(id int, eventType string) (connectors.Connector, error) {
	issuersConnections := c.ListOpenConnections(id, "")

	var result connectors.Connector
	for _, connection := range issuersConnections {
		if connection.GetEventType() == eventType {
			result = connection
		}
	}

	if result == nil {
		return nil, errors.New("connection not found with given event type")
	}

	return result, nil
}
