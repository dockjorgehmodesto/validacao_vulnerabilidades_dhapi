package connectors

import (
	"context"
	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	"errors"
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"sync"
	"testing"
)

type BillingInterceptor struct {
}

func (b BillingInterceptor) HttpInterceptor(next http.Handler) http.Handler {
	panic("implement me")
}

func (b BillingInterceptor) ServicebusInterceptor(issuerId int, eventType string, connectionInfo billing.ServiceBusConnectionInfo, callbackAfterLogging func(ctx context.Context, msg *servicebus.Message) error) servicebus.HandlerFunc {
	return callbackAfterLogging
}

func TestConnection_Disconnect(t *testing.T) {
	var cancelFunc context.CancelFunc
	var onClose func(key string)

	connector := serviceBusConnector{
		mutex: sync.Mutex{},
	}

	cancelFunc = func() {}
	onClose = func(key string) { return }

	t.Run("close only subscription", func(t *testing.T) {
		serviceBusSubscription := &MockServiceBusSubscription{}
		serviceBusSubscription.On("Close", mock.Anything).Return(nil)
		connector.subscription = serviceBusSubscription

		err := connector.Disconnect()
		assert.Nil(t, err)
	})
	t.Run("close only context", func(t *testing.T) {
		connector.contextCancelFunc = cancelFunc
		err := connector.Disconnect()
		assert.Nil(t, err)
	})
	t.Run("execute only onClose method", func(t *testing.T) {
		connector.OnClose = onClose
		err := connector.Disconnect()
		assert.Nil(t, err)
	})

	t.Run("error on closing subscription", func(t *testing.T) {
		serviceBusSubscription := &MockServiceBusSubscription{}
		serviceBusSubscription.On("Close", mock.Anything).Return(errors.New("mock"))
		connector.subscription = serviceBusSubscription
		err := connector.Disconnect()
		assert.NotNil(t, err)
	})

	t.Run("error on closing subscription with detach error", func(t *testing.T) {
		serviceBusSubscription := &MockServiceBusSubscription{}
		serviceBusSubscription.On("Close", mock.Anything).Return(&amqp.DetachError{})
		connector.subscription = serviceBusSubscription
		err := connector.Disconnect()
		assert.Nil(t, err)
	})
}
