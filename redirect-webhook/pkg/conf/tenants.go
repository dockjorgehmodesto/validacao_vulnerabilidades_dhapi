package conf

import (
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers/cache"
)

func ConfigMainTenant(cacheManager cache.MarshalerInterface) {
	tenants.TenantService = tenants.NewService(
		tenants.NewCacheRepository(tenants.NewDatabaseRepository(), cacheManager),
		tenants.NewMigrationsManager(),
	)

	tenants.TenantService.StartupMainTenant()
}

func ConfigIssuerTenants(tenantsAreUp chan bool) {
	if tenants.TenantService != nil {
		tenants.TenantService.StartupIssuerTenants()
		tenantsAreUp <- true
	}
}
