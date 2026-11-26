package encryption

import (
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/astaxie/beego/validation"
	"time"
)

type KeyVaultCreate struct {
	IssuerPublicKey         string    `json:"publicKey" valid:"Required"`
	IssuerKeyExpirationDate time.Time `json:"expirationDate" valid:"Required"`
}

func (k KeyVaultCreate) ConvertToEntity() KeyVault {
	return KeyVault{
		IssuerPublicKey:         []byte(k.IssuerPublicKey),
		IssuerKeyExpirationDate: k.IssuerKeyExpirationDate,
		Enabled:                 true,
	}
}

type KeyVaultResponse struct {
	ID                      string    `json:"id"`
	Enabled                 bool      `json:"enabled"`
	IssuerKeyExpirationDate time.Time `json:"issuerKeyExpirationDate"`
	AESKeyExpirationDate    time.Time `json:"aesKeyExpirationDate"`

	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updateAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
	IssuerId        int        `json:"issuerId"`
	IssuerPublicKey string     `json:"issuerPublicKey"`
}

func ConvertEntityToResponse(vault KeyVault) KeyVaultResponse {
	return KeyVaultResponse{
		ID:                      vault.ID,
		Enabled:                 vault.Enabled,
		IssuerKeyExpirationDate: vault.IssuerKeyExpirationDate,
		AESKeyExpirationDate:    vault.AESKeyExpirationDate,
		CreatedAt:               vault.CreatedAt,
		UpdatedAt:               vault.UpdatedAt,
		DeletedAt:               vault.DeletedAt,
		IssuerId:                vault.IssuerId,
		IssuerPublicKey:         string(vault.IssuerPublicKey),
	}
}

type KeyVaultResponsePage struct {
	helpers.PageResponse
	Records []KeyVaultResponse `json:"records"`
}

func ConvertPagination(paginator helpers.Paginator) KeyVaultResponsePage {
	slice := make([]KeyVaultResponse, 0)
	records := paginator.Records.(*[]KeyVault)

	for _, rec := range *records {
		slice = append(slice, ConvertEntityToResponse(rec))
	}

	return KeyVaultResponsePage{
		PageResponse: helpers.PageResponse{
			TotalRecord: paginator.TotalRecord,
			TotalPage:   paginator.TotalPage,
			Offset:      paginator.Offset,
			Limit:       paginator.Limit,
			Page:        paginator.Page,
			PrevPage:    paginator.PrevPage,
			NextPage:    paginator.NextPage,
		},
		Records: slice,
	}
}

type KeyVaultUpdate struct {
	IssuerPublicKey         string    `json:"publicKey" valid:"Required"`
	IssuerKeyExpirationDate time.Time `json:"expirationDate" valid:"Required"`
	Enabled                 *bool     `json:"enabled" valid:"Required"`
}

func (k KeyVaultUpdate) Valid(v *validation.Validation) {
	if k.Enabled == nil {
		_ = v.SetError("enabled", "is required")
	}
}

func (k KeyVaultUpdate) ConvertToEntity(id string) KeyVault {
	return KeyVault{
		ID:                      id,
		IssuerPublicKey:         []byte(k.IssuerPublicKey),
		IssuerKeyExpirationDate: k.IssuerKeyExpirationDate,
		Enabled:                 *k.Enabled,
	}
}

type EncryptAesResponse struct {
	EncryptedAes string `json:"encryptedAes"`
}

type KeyVaultRequest struct {
	IssuerId *int  `json:"issuerId,string" model:"issuer_id,omitempty"`
	Enabled  *bool `json:"enabled,string" model:"enabled,omitempty"`
}

func (r KeyVaultRequest) ToMap() map[string]interface{} {
	return helpers.ParseStructToMap(r)
}

type KeyVaultReplicate struct {
	IssuerId int `json:"issuerId" valid:"Required"`
}
