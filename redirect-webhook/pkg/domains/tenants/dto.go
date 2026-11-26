package tenants

import (
	"time"
)

type TenantPersist struct {
	IssuerID uint   `json:"issuer_id" valid:"Required"`
	Dialect  string `json:"dialect" valid:"Required"`
	Uri      string `json:"uri" valid:"Required"`
	User     string `json:"user" valid:"Required"`
	Password string `json:"password" valid:"Required"`
	Type     string `json:"type" valid:"Required" description:"Type of permission (READ, WRITE, READ_AND_WRITE)"`
}

type TenantUpdate struct {
	IssuerID uint   `json:"issuer_id" valid:"Required"`
	Dialect  string `json:"dialect" valid:"Required"`
	Uri      string `json:"uri" valid:"Required"`
	User     string `json:"user" valid:"Required"`
	Password string `json:"password" valid:"Required"`
	Type     string `json:"type" valid:"Required" description:"Type of permission (READ, WRITE, READ_AND_WRITE)"`
}

type TenantReplicate struct {
	IssuerID int `json:"issuerId" valid:"Required"`
}

func (t TenantUpdate) ConvertToEntity(id string) Tenant {
	return Tenant{
		ID:       id,
		IssuerID: t.IssuerID,
		Dialect:  t.Dialect,
		Uri:      t.Uri,
		User:     t.User,
		Password: t.Password,
		Type:     t.Type,
	}
}

type TenantPersistDto struct {
	IssuerID uint   `json:"issuer_id" valid:"Required"`
	Dialect  string `json:"dialect" valid:"Required"`
	Uri      string `json:"uri" valid:"Required"`
	User     string `json:"user" valid:"Required"`
	Password string `json:"password" valid:"Required"`
	Type     string `json:"type" valid:"Required" description:"Type of permission (READ, WRITE, READ_AND_WRITE)"`
}

func (t TenantPersist) ToEntity() Tenant {
	return Tenant{
		IssuerID: t.IssuerID,
		Dialect:  t.Dialect,
		Uri:      t.Uri,
		User:     t.User,
		Password: t.Password,
		Type:     t.Type,
	}
}

func (t TenantPersistDto) ToEntity() Tenant {
	return Tenant{
		IssuerID: t.IssuerID,
		Dialect:  t.Dialect,
		Uri:      t.Uri,
		User:     t.User,
		Password: t.Password,
		Type:     t.Type,
	}
}

type TenantResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IssuerID  uint      `json:"issuer_id"`
	Dialect   string    `json:"dialect"`
	Uri       string    `json:"uri"`
	User      string    `json:"user"`
	Type      string    `json:"type"`
}

func ConvertEntityToResponse(tenant Tenant) TenantResponse {
	return TenantResponse{
		ID:        tenant.ID,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
		IssuerID:  tenant.IssuerID,
		Dialect:   tenant.Dialect,
		Uri:       tenant.Uri,
		User:      tenant.User,
		Type:      tenant.Type,
	}
}
