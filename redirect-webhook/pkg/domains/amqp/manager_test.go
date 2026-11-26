package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting/connectors"
	"conductor.com.br/pier/hooks/pkg/helpers/constants"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
)

func TestConnectionManagerImpl_FindByIssuerAndEventType(t *testing.T) {
	t.Run("connection not found", func(t *testing.T) {
		connectionManager := NewConnectionManager()
		eventType, err := connectionManager.FindByIssuerAndEventType(1, "")
		assert.NotNil(t, err)
		assert.Nil(t, eventType)
		assert.Contains(t, err.Error(), "connection not found with given event type")
	})

	t.Run("finds a result", func(t *testing.T) {
		connectionManager := NewConnectionManager()
		mock := &ConnectorMock{}
		issuerId := rand.Int()
		eventType := random.String(15)
		mock.On("GetId").Return(strconv.Itoa(issuerId))
		mock.On("GetIssuerId").Return(issuerId)
		mock.On("GetEventType").Return(eventType)

		connectionManager.Connect(mock)
		connection, err := connectionManager.FindByIssuerAndEventType(issuerId, eventType)
		assert.Nil(t, err)
		assert.NotNil(t, connection)
	})
}

func TestConnectionManagerImpl_filterConnectionsByIssuerId(t *testing.T) {
	type fields struct {
		openConnections map[string]connectors.Connector
		mutex           *sync.RWMutex
	}
	type args struct {
		input    []connectors.Connector
		issuerId int
	}

	connection1, _ := connectors.CreateConnector(connecting.Connection{ID: "1", IssuerId: 1, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection2, _ := connectors.CreateConnector(connecting.Connection{ID: "2", IssuerId: 1, EventType: constants.EventTypeAddress, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection3, _ := connectors.CreateConnector(connecting.Connection{ID: "3", IssuerId: 1, EventType: constants.EventTypeStatement, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection4, _ := connectors.CreateConnector(connecting.Connection{ID: "4", IssuerId: 2, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection5, _ := connectors.CreateConnector(connecting.Connection{ID: "5", IssuerId: 2, EventType: constants.EventTypeCard, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection6, _ := connectors.CreateConnector(connecting.Connection{ID: "6", IssuerId: 3, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection7, _ := connectors.CreateConnector(connecting.Connection{ID: "7", IssuerId: 4, EventType: constants.EventTypePurchase, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connectionsSample := []connectors.Connector{
		connection1,
		connection2,
		connection3,
		connection4,
		connection5,
		connection6,
		connection7,
	}

	fieldsSample := fields{
		make(map[string]connectors.Connector),
		&sync.RWMutex{},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []connectors.Connector
	}{
		{"filtering by issuerId 1", fieldsSample, args{connectionsSample, 1}, []connectors.Connector{connection1, connection2, connection3}},
		{"filtering by issuerId 2", fieldsSample, args{connectionsSample, 2}, []connectors.Connector{connection4, connection5}},
		{"filtering by issuerId 3", fieldsSample, args{connectionsSample, 3}, []connectors.Connector{connection6}},
		{"filtering by issuerId 4", fieldsSample, args{connectionsSample, 4}, []connectors.Connector{connection7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnectionManagerImpl{
				openConnections: tt.fields.openConnections,
				mutex:           tt.fields.mutex,
			}
			if got := c.filterConnectionsByIssuerId(tt.args.input, tt.args.issuerId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterConnectionsByIssuerId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnectionManagerImpl_filterConnectionsByEventType(t *testing.T) {

	type fields struct {
		openConnections map[string]connectors.Connector
		mutex           *sync.RWMutex
	}
	type args struct {
		input     []connectors.Connector
		eventType string
	}

	connection1, _ := connectors.CreateConnector(connecting.Connection{ID: "1", IssuerId: 1, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection2, _ := connectors.CreateConnector(connecting.Connection{ID: "2", IssuerId: 1, EventType: constants.EventTypeAddress, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection3, _ := connectors.CreateConnector(connecting.Connection{ID: "3", IssuerId: 1, EventType: constants.EventTypeStatement, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection4, _ := connectors.CreateConnector(connecting.Connection{ID: "4", IssuerId: 2, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection5, _ := connectors.CreateConnector(connecting.Connection{ID: "5", IssuerId: 2, EventType: constants.EventTypeCard, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection6, _ := connectors.CreateConnector(connecting.Connection{ID: "6", IssuerId: 3, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection7, _ := connectors.CreateConnector(connecting.Connection{ID: "7", IssuerId: 4, EventType: constants.EventTypePurchase, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connectionsSample := []connectors.Connector{
		connection1,
		connection2,
		connection3,
		connection4,
		connection5,
		connection6,
		connection7,
	}

	fieldsSample := fields{
		make(map[string]connectors.Connector),
		&sync.RWMutex{},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []connectors.Connector
	}{
		{"filtering by event account", fieldsSample, args{connectionsSample, constants.EventTypeAccount}, []connectors.Connector{connection1, connection4, connection6}},
		{"filtering by event address", fieldsSample, args{connectionsSample, constants.EventTypeAddress}, []connectors.Connector{connection2}},
		{"filtering by event statement", fieldsSample, args{connectionsSample, constants.EventTypeStatement}, []connectors.Connector{connection3}},
		{"filtering by event account", fieldsSample, args{connectionsSample, constants.EventTypeAccount}, []connectors.Connector{connection1, connection4, connection6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnectionManagerImpl{
				openConnections: tt.fields.openConnections,
				mutex:           tt.fields.mutex,
			}
			if got := c.filterConnectionsByEventType(tt.args.input, tt.args.eventType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterConnectionsByEventType() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestConnectionManagerImpl_ListOpenConnections(t *testing.T) {

	type fields struct {
		openConnections map[string]connectors.Connector
		mutex           *sync.RWMutex
	}
	type args struct {
		issuerId  int
		eventType string
	}

	connection1, _ := connectors.CreateConnector(connecting.Connection{ID: "1", IssuerId: 1, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection2, _ := connectors.CreateConnector(connecting.Connection{ID: "2", IssuerId: 1, EventType: constants.EventTypeAddress, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection3, _ := connectors.CreateConnector(connecting.Connection{ID: "3", IssuerId: 1, EventType: constants.EventTypeStatement, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection4, _ := connectors.CreateConnector(connecting.Connection{ID: "4", IssuerId: 2, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection5, _ := connectors.CreateConnector(connecting.Connection{ID: "5", IssuerId: 2, EventType: constants.EventTypeCard, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection6, _ := connectors.CreateConnector(connecting.Connection{ID: "6", IssuerId: 3, EventType: constants.EventTypeAccount, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connection7, _ := connectors.CreateConnector(connecting.Connection{ID: "7", IssuerId: 4, EventType: constants.EventTypePurchase, Broker: connecting.AMQPBrokerRabbitMQ}, nil, nil)
	connectionsSample := []connectors.Connector{
		connection1,
		connection2,
		connection3,
		connection4,
		connection5,
		connection6,
		connection7,
	}

	fieldsSample := fields{
		map[string]connectors.Connector{
			"1": connection1,
			"2": connection2,
			"3": connection3,
			"4": connection4,
			"5": connection5,
			"6": connection6,
			"7": connection7,
		},
		&sync.RWMutex{},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []connectors.Connector
	}{
		{"listing all results", fieldsSample, args{0, ""}, connectionsSample},
		{"filtering by issuer id 1", fieldsSample, args{1, ""}, []connectors.Connector{connection1, connection2, connection3}},
		{"filtering by issuer id 2", fieldsSample, args{2, ""}, []connectors.Connector{connection4, connection5}},
		{"filtering by eventType account", fieldsSample, args{0, constants.EventTypeAccount}, []connectors.Connector{connection4, connection6, connection1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnectionManagerImpl{
				openConnections: tt.fields.openConnections,
				mutex:           tt.fields.mutex,
			}
			if got := c.ListOpenConnections(tt.args.issuerId, tt.args.eventType); !reflect.DeepEqual(sortConnectionsById(got), sortConnectionsById(tt.want)) {
				t.Errorf("ListOpenConnections() = len(%v), want len(%v)", len(got), len(tt.want))
				t.Errorf("ListOpenConnections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortConnectionsById(slice []connectors.Connector) []connectors.Connector {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].GetId() < slice[j].GetId()
	})
	return slice
}

func TestConnectionManagerImpl_Connect(t *testing.T) {
	type fields struct {
		openConnections map[string]connectors.Connector
		mutex           *sync.RWMutex
	}
	type args struct {
		connection connectors.Connector
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"success", fields{openConnections: make(map[string]connectors.Connector), mutex: &sync.RWMutex{}}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnectionManagerImpl{
				openConnections: tt.fields.openConnections,
				mutex:           tt.fields.mutex,
			}

			connection := &ConnectorMock{}
			connection.On("GetId").Return(random.String(15))

			c.Connect(connection)
			//assert.NotPanics(t, func() {
			//})
		})
	}
}

func TestConnectionManagerImpl_DeleteByKey(t *testing.T) {

	connectionMock := &ConnectorMock{}
	connectionMock.On("GetId").Return("1")
	mapMock := map[string]connectors.Connector{
		"1": connectionMock,
	}

	type fields struct {
		openConnections map[string]connectors.Connector
		mutex           *sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"passing invalid id", fields{make(map[string]connectors.Connector), &sync.RWMutex{}}, args{"invalid"}},
		{"passing valid id with success", fields{mapMock, &sync.RWMutex{}}, args{"1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConnectionManagerImpl{
				openConnections: tt.fields.openConnections,
				mutex:           tt.fields.mutex,
			}
			assert.NotPanics(t, func() {
				c.DeleteByKey(tt.args.key)
			})
		})
	}
}
