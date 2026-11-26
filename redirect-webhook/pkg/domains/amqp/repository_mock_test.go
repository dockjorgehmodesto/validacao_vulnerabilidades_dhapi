package amqp

import (
	"conductor.com.br/pier/hooks/pkg/domains/amqp/connecting"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (f *RepositoryMock) FindById(key *tenants.TenantKey, id string) (connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) FindByEventType(key *tenants.TenantKey, eventType string) (connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, example map[string]interface{}) (helpers.Paginator, error) {
	panic("implement me")
}

func (f *RepositoryMock) Save(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) Update(key *tenants.TenantKey, connection connecting.Connection) (connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) UpdateStatus(key *tenants.TenantKey, connection connecting.Connection) error {
	return f.expectedCallsMock("UpdateStatus").Get(0).(error)
}

func (f *RepositoryMock) PartialUpdate(key *tenants.TenantKey, id string, values map[string]interface{}) (connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) ListEnabled(key *tenants.TenantKey) ([]connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) ListEnabledByIssuer(key *tenants.TenantKey) ([]connecting.Connection, error) {
	panic("implement me")
}

func (f *RepositoryMock) Delete(key *tenants.TenantKey, id string) error {
	return f.expectedCallsMock("Delete").Get(0).(error)
}

func (f *RepositoryMock) expectedCallsMock(method string) mock.Arguments {
	var mockArguments mock.Arguments

	for i := 0; i < len(f.ExpectedCalls); i++ {
		if f.ExpectedCalls[i].Method == method {
			mockArguments = f.ExpectedCalls[i].ReturnArguments
		}
	}
	return mockArguments
}
