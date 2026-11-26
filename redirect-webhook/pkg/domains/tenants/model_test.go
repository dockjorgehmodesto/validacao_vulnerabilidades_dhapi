package tenants

import (
	"github.com/gofrs/uuid"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestModel_GetTenantKey(t *testing.T) {
	t.Run("from empty tenant", func(t *testing.T) {
		emptyTenant := Tenant{}
		tenantKey := emptyTenant.GetTenantKey()
		assert.Equal(t, "0_", tenantKey.Key)
	})

	t.Run("from tenant with empty type and issuer id", func(t *testing.T) {
		randString := random.String(10)
		tenant := Tenant{
			ID: randString,
		}
		tenantKey := tenant.GetTenantKey()
		assert.Equal(t, "0_", tenantKey.Key)
	})

	t.Run("from tenant with empty type", func(t *testing.T) {
		randString := random.String(10)
		issuerId := uint(rand.Uint32())
		tenant := Tenant{
			ID:       randString,
			IssuerID: issuerId,
		}
		tenantKey := tenant.GetTenantKey()
		assert.Equal(t, strconv.Itoa(int(issuerId))+"_", tenantKey.Key)
	})

	t.Run("from tenant with all data", func(t *testing.T) {
		randString := random.String(10)
		issuerId := uint(rand.Uint32())
		tenant := Tenant{
			ID:       randString,
			IssuerID: issuerId,
			Type:     TenantTypeRead,
		}
		tenantKey := tenant.GetTenantKey()
		assert.Equal(t, strconv.Itoa(int(issuerId))+"_"+TenantTypeRead, tenantKey.Key)
	})
}

func TestModel_FormatUriConnection(t *testing.T) {
	t.Run("with empty tenant", func(t *testing.T) {
		tenant := Tenant{}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, " user= password=", uriConnection)
	})

	t.Run("with empty user and password", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		tenant := Tenant{
			Uri: uri,
		}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, uri+" user= password=", uriConnection)
	})

	t.Run("with empty user", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		pass := random.String(15)
		tenant := Tenant{
			Uri:      uri,
			Password: pass,
		}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, uri+" user= password="+pass, uriConnection)
	})

	t.Run("with empty password", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		user := random.String(15)
		tenant := Tenant{
			Uri:  uri,
			User: user,
		}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, uri+" user="+user+" password=", uriConnection)
	})

	t.Run("with empty uri", func(t *testing.T) {
		user := random.String(15)
		pass := random.String(15)
		tenant := Tenant{
			User:     user,
			Password: pass,
		}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, " user="+user+" password="+pass, uriConnection)
	})

	t.Run("valid scenario", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		user := random.String(15)
		pass := random.String(15)
		tenant := Tenant{
			Uri:      uri,
			User:     user,
			Password: pass,
		}
		uriConnection := tenant.FormatUriConnection()
		assert.Equal(t, uri+" user="+user+" password="+pass, uriConnection)
	})
}

func Test_CreateTenantKeyString(t *testing.T) {
	t.Run("from random tenant with identifier", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		user := random.String(15)
		pass := random.String(15)
		tenant := Tenant{
			ID:       random.String(15),
			Uri:      uri,
			User:     user,
			Password: pass,
		}

		tenantId := CreateTenantKeyString(&tenant)
		assert.NotNil(t, tenantId)
		assert.NotEmpty(t, tenantId)
		assert.NotEqual(t, tenantId, TENANT_PRINCIPAL_ID)
	})

	t.Run("from main tenant (without identifier)", func(t *testing.T) {
		uri := "host=myconnection port=1234 ssl=disable"
		user := random.String(15)
		pass := random.String(15)
		tenant := Tenant{
			Uri:      uri,
			User:     user,
			Password: pass,
		}

		tenantId := CreateTenantKeyString(&tenant)
		assert.NotNil(t, tenantId)
		assert.NotEmpty(t, tenantId)
		assert.Equal(t, tenantId, TENANT_PRINCIPAL_ID)
	})
}

func Test_CreateTenantKeyFromIdentifier(t *testing.T) {
	t.Run("with empty string", func(t *testing.T) {
		assert.Equal(t, "", CreateTenantKeyFromIdentifier("").Key)
	})

	t.Run("with uuid", func(t *testing.T) {
		uuids, _ := uuid.NewV4()
		id := uuids.String()
		tenantKey := CreateTenantKeyFromIdentifier(id)
		assert.NotEmpty(t, tenantKey.Key)
		assert.Equal(t, id, tenantKey.Key)
	})
}
