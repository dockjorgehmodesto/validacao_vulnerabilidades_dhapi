package encryption

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"context"
	"time"
)

func NewService(repository Repository, tenantService tenants.Service) Service {
	return ServiceImpl{repository: repository, tenantService: tenantService}
}

type Service interface {
	Create(key *tenants.TenantKey, vault KeyVault) KeyVault
	Update(key *tenants.TenantKey, keyVault KeyVault) KeyVault
	FindById(key *tenants.TenantKey, id string) KeyVault
	FindFirst(key *tenants.TenantKey) (KeyVault, error)
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator
	EncryptAes(key *tenants.TenantKey, keyVaultId string) string
	FindMostRecentAndActive(key *tenants.TenantKey) (KeyVault, error)
	FindMostRecentAndActiveByIssuerId(ctx context.Context, key *tenants.TenantKey, issuerId int) (KeyVault, error)
	CreateMockKeyVault(publicKey string) KeyVault
	Replicate(tenant *tenants.TenantKey, keyVaultId string, issuerId int) KeyVault
}

type ServiceImpl struct {
	repository    Repository
	tenantService tenants.Service
}

func (s ServiceImpl) FindFirst(key *tenants.TenantKey) (KeyVault, error) {
	return s.repository.FindFirst(key)
}

func (s ServiceImpl) EncryptAes(key *tenants.TenantKey, keyVaultId string) string {
	keyVault := s.FindById(key, keyVaultId)

	Error{}.HandleExpiredPublicKey(keyVault.IssuerKeyHasExpired())
	Error{}.HandleDisabledVault(keyVault.Enabled)

	result, e := keyVault.EncryptAesWithPublicKey()
	Error{}.HandleAesEncryption(e)
	return result
}

func (s ServiceImpl) Create(key *tenants.TenantKey, vault KeyVault) KeyVault {

	s.validateAndSetupKeyVault(&vault)
	savedKeyVault, e := s.repository.Save(key, vault)
	Error{}.HandleDatabaseError(e)
	return savedKeyVault
}

func (s ServiceImpl) validateAndSetupKeyVault(vault *KeyVault) {
	Error{}.HandleInvalidPublicKey(vault.ValidatePublicKey())
	Error{}.HandleExpiredPublicKey(vault.IssuerKeyHasExpired())

	err := vault.CreateAesKey()
	Error{}.HandleAesKey(err)

	vault.SetInitialStatusAndExpiration()
}

func (s ServiceImpl) validateAndSetupKeyVaultMock(vault *KeyVault) {
	Error{}.HandleInvalidPublicKey(vault.ValidatePublicKey())
	Error{}.HandleExpiredPublicKey(vault.IssuerKeyHasExpired())

	err := vault.CreateAesKeyMock()
	Error{}.HandleAesKey(err)

	vault.SetInitialStatusAndExpiration()
}

func (s ServiceImpl) Update(key *tenants.TenantKey, keyVault KeyVault) KeyVault {
	keyVaultToUpdate := s.FindById(key, keyVault.ID)

	// TODO Change this to work with a map
	ParseUpdateData(&keyVaultToUpdate, keyVault)
	updatedKeyVault, e := s.repository.Update(key, keyVault)
	Error{}.HandleDatabaseError(e)
	return updatedKeyVault
}

// TODO Change this to work with a map
func ParseUpdateData(dst *KeyVault, src KeyVault) {
	dst.Enabled = src.Enabled
	dst.IssuerPublicKey = src.IssuerPublicKey
	dst.IssuerKeyExpirationDate = src.IssuerKeyExpirationDate
}

func (s ServiceImpl) FindById(key *tenants.TenantKey, id string) KeyVault {
	keyVault, e := s.repository.FindById(key, id)
	Error{}.HandleSearch(e)
	return keyVault
}

func (s ServiceImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string, queryExample map[string]interface{}) helpers.Paginator {
	paginator, e := s.repository.FindAllWithPagination(key, page, limit, sort, queryExample)
	Error{}.HandlePagination(paginator, e)
	return paginator
}

func (s ServiceImpl) FindMostRecentAndActive(key *tenants.TenantKey) (KeyVault, error) {
	return s.repository.FindMostRecentAndActive(key)
}

func (s ServiceImpl) FindMostRecentAndActiveByIssuerId(ctx context.Context, key *tenants.TenantKey, issuerId int) (KeyVault, error) {
	return s.repository.FindMostRecentAndActiveByIssuerId(ctx, key, issuerId)
}

func (s ServiceImpl) CreateMockKeyVault(publicKey string) KeyVault {
	keyVault := KeyVault{
		IssuerPublicKey:         []byte(publicKey),
		IssuerKeyExpirationDate: time.Now().Add(time.Hour * 48),
	}

	s.validateAndSetupKeyVaultMock(&keyVault)
	return keyVault
}

func (s ServiceImpl) Replicate(tenant *tenants.TenantKey, keyVaultId string, issuerId int) KeyVault {
	// Validating if target issuer ID is valid
	_, err := s.tenantService.FindOneByIssuerId(issuerId)
	Error{}.HandleInvalidTargetIssuer(err)

	// Validating if source key vault exists
	keyVault := s.FindById(tenant, keyVaultId)

	keyVault.ID = ""
	keyVault.IssuerId = issuerId
	keyVault.CreatedAt = time.Now()
	keyVault.UpdatedAt = time.Now()
	savedKeyVault, e := s.repository.Save(tenant, keyVault)
	Error{}.HandleDatabaseError(e)
	return savedKeyVault
}
