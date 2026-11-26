package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"github.com/Azure/azure-service-bus-go"
	"github.com/stretchr/testify/mock"
)

type ConnectorMock struct {
	mock.Mock
}

func (c *ConnectorMock) GetProcessorType() string {
	panic("implement me")
}

func (c *ConnectorMock) GetConnection() *connecting.Connection {
	panic("implement me")
}

func (f *ConnectorMock) expectedCallsMock(method string) mock.Arguments {
	var mockArguments mock.Arguments

	for i := 0; i < len(f.ExpectedCalls); i++ {
		if f.ExpectedCalls[i].Method == method {
			mockArguments = f.ExpectedCalls[i].ReturnArguments
		}
	}
	return mockArguments
}

func (c *ConnectorMock) GetConnectionStringMap() map[string]string {
	return c.expectedCallsMock("GetConnectionStringMap").Get(0).(map[string]string)
}

func (c *ConnectorMock) GetMaskedConnectionString() string {
	return c.expectedCallsMock("GetMaskedConnectionString").Get(0).(string)
}

func (c *ConnectorMock) GetNamespaceWithConnectionString() servicebus.NamespaceOption {
	return c.expectedCallsMock("GetNamespaceWithConnectionString").Get(0).(servicebus.NamespaceOption)
}

func (c *ConnectorMock) Test() error {
	return c.expectedCallsMock("Test").Get(0).(error)
}

func (c *ConnectorMock) CountDeadLetterMessages() (int, error) {
	args := c.expectedCallsMock("CountDeadLetterMessages")
	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return args.Get(0).(int), err
}

func (c ConnectorMock) ResendDeadLetteredMessages() (int, error) {
	args := c.expectedCallsMock("ResendDeadLetteredMessages")
	var err error
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return args.Get(0).(int), err
}

func (c *ConnectorMock) Listen() {
}

func (c *ConnectorMock) Disconnect() error {
	if c.expectedCallsMock("Disconnect").Get(0) == nil {
		return nil
	}
	if val, ok := c.expectedCallsMock("Disconnect").Get(0).(error); ok {
		return val
	}
	return nil
}

func (c *ConnectorMock) GetEventType() string {
	return c.expectedCallsMock("GetEventType").Get(0).(string)
}

func (c *ConnectorMock) GetId() string {
	return c.expectedCallsMock("GetId").Get(0).(string)
}

func (c *ConnectorMock) GetIssuerId() int {
	return c.expectedCallsMock("GetIssuerId").Get(0).(int)
}

func (c *ConnectorMock) SetOnClose(onclose func(key string)) {
}
