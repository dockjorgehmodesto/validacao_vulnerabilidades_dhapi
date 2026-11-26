package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/domains/webhooks"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/stretchr/testify/mock"
)

type EncryptionServiceMock struct {
	mock.Mock
}

func (s *EncryptionServiceMock) Replicate(tenant *tenants.TenantKey, keyVaultId string, issuerId int) encryption.KeyVault {
	panic("implement me")
}

func NewEncryptionMockService() *EncryptionServiceMock {
	return &EncryptionServiceMock{}
}

func (s *EncryptionServiceMock) Create(key *tenants.TenantKey, vault encryption.KeyVault) encryption.KeyVault {
	panic("implement me")
}

func (s *EncryptionServiceMock) Update(key *tenants.TenantKey, keyVault encryption.KeyVault) encryption.KeyVault {
	panic("implement me")
}

func (s *EncryptionServiceMock) FindById(key *tenants.TenantKey, id string) encryption.KeyVault {
	panic("implement me")
}

func (s *EncryptionServiceMock) FindFirst(key *tenants.TenantKey) (encryption.KeyVault, error) {
	panic("implement me")
}

func (s *EncryptionServiceMock) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator {
	panic("implement me")
}

func (s *EncryptionServiceMock) EncryptAes(key *tenants.TenantKey, keyVaultId string) string {
	panic("implement me")
}

func (s *EncryptionServiceMock) FindMostRecentAndActive(key *tenants.TenantKey) (encryption.KeyVault, error) {
	args := s.expectedCallsMock("FindMostRecentAndActive")
	var secondArg error
	if args.Get(1) != nil {
		secondArg = args.Get(1).(error)
	}
	return args.Get(0).(encryption.KeyVault), secondArg
}

func (s *EncryptionServiceMock) CreateMockKeyVault(publicKey string) encryption.KeyVault {
	panic("implement me")
}

func (f *EncryptionServiceMock) expectedCallsMock(method string) mock.Arguments {
	var mockArguments mock.Arguments

	for i := 0; i < len(f.ExpectedCalls); i++ {
		if f.ExpectedCalls[i].Method == method {
			mockArguments = f.ExpectedCalls[i].ReturnArguments
		}
	}
	return mockArguments
}

type WebhooksServiceMock struct {
	mock.Mock
}

func NewWebhookMockService() *WebhooksServiceMock {
	return &WebhooksServiceMock{}
}

func (s *WebhooksServiceMock) Create(key *tenants.TenantKey, hook webhooks.WebHook) webhooks.WebHook {
	panic("implement me")
}

func (s *WebhooksServiceMock) Update(key *tenants.TenantKey, id string, updateMap map[string]interface{}) webhooks.WebHook {
	panic("implement me")
}

func (s *WebhooksServiceMock) FindByIdWithPreload(key *tenants.TenantKey, id string) webhooks.WebHook {
	return s.expectedCallsMock("FindByIdWithPreload").Get(0).(webhooks.WebHook)
}

func (s *WebhooksServiceMock) FetchWithPreload(key *tenants.TenantKey, id string) (webhooks.WebHook, error) {
	panic("implement me")
}

func (s *WebhooksServiceMock) ValidateIfTypeIsAvailable(key *tenants.TenantKey, eventType string) {
	panic("implement me")
}

func (s *WebhooksServiceMock) FindByTypeWithPreload(key *tenants.TenantKey, eventType string) (webhooks.WebHook, error) {
	args := s.expectedCallsMock("FindByTypeWithPreload")
	var secondArg error
	if args.Get(1) != nil {
		secondArg = args.Get(1).(error)
	}
	return args.Get(0).(webhooks.WebHook), secondArg
}

func (s *WebhooksServiceMock) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string) helpers.Paginator {
	panic("implement me")
}

func (s *WebhooksServiceMock) Trigger(webhook webhooks.WebHook, payload string, forceMockPayload bool, eventType string, publicKey string) (statusCode int, responseBody string, encryptedAesKey string) {
	panic("implement me")
}

func (s *WebhooksServiceMock) broadcastFlushConnections(tenantKey *tenants.TenantKey) {
	panic("implement me")
}

func (s *WebhooksServiceMock) Delete(tenant *tenants.TenantKey, id string) {
	panic("implement me")
}

func (f *WebhooksServiceMock) expectedCallsMock(method string) mock.Arguments {
	var mockArguments mock.Arguments

	for i := 0; i < len(f.ExpectedCalls); i++ {
		if f.ExpectedCalls[i].Method == method {
			mockArguments = f.ExpectedCalls[i].ReturnArguments
		}
	}
	return mockArguments
}
