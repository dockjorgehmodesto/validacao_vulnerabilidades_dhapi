package main

import (
	_ "net/http/pprof"
	"time"

	httpConf "conductor.com.br/pier/hooks/cmd/http/conf"
	"conductor.com.br/pier/hooks/pkg/conf"
	"conductor.com.br/pier/hooks/pkg/domains"
	"conductor.com.br/pier/hooks/pkg/infraestructure/observability"
	// "dev.azure.com/conductortech/vnext/_git/core.git"

	"conductor.com.br/pier/hooks/pkg/helpers/logging"
)

// @title DataHub API
// @version 1.0
// @description API documentation for PIER Hooks microservice
// @termsOfService http://swagger.io/terms/

// @contact.name conductor.com.br
// @contact.url http://conductor.com.br/
// @contact.email contato@conductor.com.br

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /v1

// @securityDefinitions.apikey issuerIdAuthentication
// @in header
// @name issuer_id
func main() {
	timeStarted := time.Now()

	conf.ConfigVersion()
	conf.SetupBilling()

	manager := conf.CreateCacheManager()
	marshaler := conf.ParseCacheManagerToMarshaler(manager)

	tenantsAreUp := make(chan bool)
	conf.ConfigMainTenant(marshaler)
	go conf.ConfigIssuerTenants(tenantsAreUp)

	observerMetrics := observability.NewDatadogObserver()
	defer observerMetrics.Close()

	servicesContainer := domains.Container(domains.CacheContext{Manager: manager, Marshaler: marshaler}, observerMetrics)
	go startAfterIssuerTenants(servicesContainer, tenantsAreUp)

	conf.ConfigRabbitMQ(servicesContainer)
	httpConf.ConfigChi(servicesContainer, timeStarted)
}

// func startDataDogProfilerAndTracer() {
// 	err := profiler.Start(
// 		profiler.WithService(viper.GetString("app.serviceName")),
// 		profiler.WithEnv(viper.GetString("app.profile")),
// 		profiler.WithVersion(global.VersionInfo.Current),
// 		profiler.WithTags(fmt.Sprintf("pod_name:%s", viper.GetString("pod.name"))),
// 	)
// 	if err != nil {
// 		core.Console().Error("datadog profiler error", zap.Error(err))
// 	}
// 	tracer.Start(
// 		tracer.WithService(viper.GetString("app.serviceName")),
// 		tracer.WithEnv(viper.GetString("app.profile")),
// 		tracer.WithRuntimeMetrics(),
// 		tracer.WithServiceVersion(global.VersionInfo.Current),
// 		tracer.WithGlobalTag("pod_name", viper.GetString("pod.name")),
// 	)
// }

func startAfterIssuerTenants(servicesContainer domains.ServiceContainer, tenantsAreUp chan bool) {
	<-tenantsAreUp
	logging.Logger.Info("[STARTUP] Listening to messages")
	conf.ListenToMessages(servicesContainer)
	logging.Logger.Info("[STARTUP] Updating eventController")
	go conf.UpdateEventController(servicesContainer)
	logging.Logger.Info("[STARTUP] Starting deliverer")

	for i := 0; i < len(servicesContainer.WebHookDeliverer); i++ {
		servicesContainer.WebHookDeliverer[i].StartListening()
	}

}
