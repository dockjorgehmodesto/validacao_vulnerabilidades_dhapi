package conf

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
)

func ListenToMessages(container domains.ServiceContainer) {
	if !global.Viper.GetBool("features.disableMessageReading") {
		tenantSlice, e := tenants.TenantService.ListAll()
		if e != nil {
			panic("could not load tenants")
		}
		go func() {
			for _, v := range tenantSlice {
				tenantKey := v.GetTenantKey()
				container.AMQP.RecheckConnections(&tenantKey)
			}
		}()
	}
}
