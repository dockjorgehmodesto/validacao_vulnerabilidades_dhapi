package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting/connectors"
	"conductor.com.br/pier/hooks/pkg/helpers/constants"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"sync"
	"testing"
)

func createConnectionMap() map[string]connectors.Connector {
	connectionMap := make(map[string]connectors.Connector)

	for i := 0; i < 100; i++ {
		identifier := strconv.Itoa(i)
		connection := &ConnectorMock{}
		connection.On("GetEventType").Return(constants.EventTypeAccount)
		connection.On("GetId").Return(identifier)

		// Issuer with id 1 will throw an error
		connection.On("GetIssuerId").Return(i % 2)
		if i%2 == 1 {
			connection.On("Disconnect").Return(errors.New("mock"))
		} else {
			connection.On("Disconnect").Return(nil)
		}
		connectionMap[identifier] = connection
	}
	return connectionMap
}

func TestServiceImpl_cleanupConnectionsForIssuer(t *testing.T) {

	type fields struct {
		connectionMap ConnectionManager
	}
	type args struct {
		issuerId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty connection map", fields{&ConnectionManagerImpl{openConnections: make(map[string]connectors.Connector), mutex: &sync.RWMutex{}}}, args{0}},
		{"search with no results", fields{&ConnectionManagerImpl{openConnections: createConnectionMap(), mutex: &sync.RWMutex{}}}, args{999}},
		{"errors when disconnecting", fields{&ConnectionManagerImpl{openConnections: createConnectionMap(), mutex: &sync.RWMutex{}}}, args{1}},
		{"success", fields{&ConnectionManagerImpl{openConnections: createConnectionMap(), mutex: &sync.RWMutex{}}}, args{2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceImpl{
				connectionMap: tt.fields.connectionMap,
			}
			assert.NotPanics(t, func() {
				s.cleanupConnectionsForIssuer(tt.args.issuerId)
			})
		})
	}
}

func TestServiceImpl_Disconnect(t *testing.T) {
	type fields struct {
		connectionMap ConnectionManager
	}
	type args struct {
		connectionId string
	}

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantedMapLength int
	}{
		{"invalid connectionId", fields{&ConnectionManagerImpl{mutex: &sync.RWMutex{}, openConnections: createConnectionMap()}}, args{""}, true, 100},
		{"error disconnecting", fields{&ConnectionManagerImpl{mutex: &sync.RWMutex{}, openConnections: createConnectionMap()}}, args{"3"}, true, 100},
		{"success", fields{&ConnectionManagerImpl{mutex: &sync.RWMutex{}, openConnections: createConnectionMap()}}, args{"2"}, false, 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ServiceImpl{
				connectionMap: tt.fields.connectionMap,
			}
			if err := s.Disconnect(tt.args.connectionId); (err != nil) != tt.wantErr {
				t.Errorf("Disconnect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if list := s.connectionMap.ListOpenConnections(0, ""); len(list) != tt.wantedMapLength {
				t.Errorf("List size is %d, should be %d", len(list), tt.wantedMapLength)
			}
		})
	}
}

func TestCanConnectBasedOnEnvironment(t *testing.T) {
	type args struct {
		connection connecting.Connection
	}
	connection := connecting.Connection{IssuerId: 1, EventType: "purchase"}
	tests := []struct {
		name  string
		args  args
		want  bool
		setup func()
	}{
		{"feature disabled", args{connection: connection}, true, func() { _ = os.Setenv("EVENTSELECTION_ENABLED", "false") }},
		{"valid issuer and event", args{connection: connection}, true, func() {
			_ = os.Setenv("EVENTSELECTION_ENABLED", "true")
			_ = os.Setenv("EVENTSELECTION_ISSUERS", "1")
			_ = os.Setenv("EVENTSELECTION_EVENTS", "purchase")
		}},
		{"valid issuer and invalid event", args{connection: connection}, false, func() {
			_ = os.Setenv("EVENTSELECTION_ENABLED", "true")
			_ = os.Setenv("EVENTSELECTION_ISSUERS", "1")
			_ = os.Setenv("EVENTSELECTION_EVENTS", "invalid")
		}},
		{"invalid issuer and invalid event", args{connection: connection}, false, func() {
			_ = os.Setenv("EVENTSELECTION_ENABLED", "true")
			_ = os.Setenv("EVENTSELECTION_ISSUERS", "2")
			_ = os.Setenv("EVENTSELECTION_EVENTS", "invalid")
		}},
		{"issuer and event list testing", args{connection: connection}, true, func() {
			_ = os.Setenv("EVENTSELECTION_ENABLED", "true")
			_ = os.Setenv("EVENTSELECTION_ISSUERS", "1 2 3 4")
			_ = os.Setenv("EVENTSELECTION_EVENTS", "purchase statement")
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := canConnectBasedOnEnvironment(tt.args.connection); got != tt.want {
				t.Errorf("canConnectBasedOnEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
