package conf

import (
	"conductor.com.br/pier/hooks/pkg/domains"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
)

func UpdateEventController(container domains.ServiceContainer) {
	key, err := container.Tenant.GetIssuerTenantKey()
	if err != nil {
		panic(err)
	}
	container.EventController.UpdateEvents(key)
	logging.Logger.Info("[EVENTCONTROLLER] Updated")
}
