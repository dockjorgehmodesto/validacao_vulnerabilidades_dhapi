package tenants

import (
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
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

func (r *RepositoryMock) Create(t TenantContract) (Tenant, error) {
	args := r.expectedCallsMock("Create")
	return args.Get(0).(Tenant), args.Error(1)
}

func (r RepositoryMock) Delete(id string) error {
	args := r.expectedCallsMock("Delete")
	return args.Error(0)
}

func (r RepositoryMock) Update(id string, t *Tenant) error {
	args := r.expectedCallsMock("Update")
	return args.Error(0)
}

func (r RepositoryMock) FindById(id string) (Tenant, error) {
	args := r.expectedCallsMock("FindById")
	return args.Get(0).(Tenant), args.Error(1)
}

func (r RepositoryMock) FindByIssuerIdAndType(issuerId int, accessType string) (TenantContract, error) {
	args := r.expectedCallsMock("FindByIssuerIdAndType")
	return args.Get(0).(TenantContract), args.Error(1)
}

func (r RepositoryMock) FindAll() ([]Tenant, error) {
	args := r.expectedCallsMock("FindAll")
	return args.Get(0).([]Tenant), args.Error(1)
}

func (r RepositoryMock) FindAllByPagination(page int, limit int, orderBy []string) ([]Tenant, error) {
	args := r.expectedCallsMock("ListAllByPagination")
	return args.Get(0).([]Tenant), args.Error(1)
}

func (r RepositoryMock) Ping(tenantKey *TenantKey) (SQLResult, error) {
	args := r.expectedCallsMock("Ping")
	return args.Get(0).(SQLResult), args.Error(1)
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}
