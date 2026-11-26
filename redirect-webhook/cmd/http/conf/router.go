package conf

import (
	// "conductor.com.br/pier/hooks/cmd/http/docs"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"conductor.com.br/pier/hooks/cmd/http/helpers"
	"conductor.com.br/pier/hooks/cmd/http/middleware"
	"conductor.com.br/pier/hooks/cmd/http/routes"
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/go-chi/chi"
	middleware2 "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gobuffalo/packr"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	// "github.com/swaggo/http-swagger"
)

type ChiRouterContract interface {
	CreateRoutes()
	StartHttpServer()
}

type ChiRouter struct {
	Router chi.Router
}

type RouterContainer struct {
	Internals      helpers.Router
	Management     helpers.Router
	AMQP           helpers.Router
	Encryption     helpers.Router
	WebHooks       helpers.Router
	Events         helpers.Router
	IssuerEvents   helpers.Router
	DeliveryQueues helpers.Router
}

func NewRouterContainer(container domains.ServiceContainer) RouterContainer {

	return RouterContainer{
		Internals: routes.Internals{
			Tenants: routes.NewTenantsRouter(container.Tenant),
		},
		Management:     routes.NewManagementRouter(container.Management),
		AMQP:           routes.NewAMQPRouter(container.AMQP, container.Tenant),
		WebHooks:       routes.NewWebhooksRouter(container.WebHook),
		Encryption:     routes.NewEncryptionRouter(container.Encryption),
		Events:         routes.NewEventTypesRouter(container.AMQP, container.IssuerEvents),
		IssuerEvents:   routes.NewIssuerEventsRouter(container.IssuerEvents),
		DeliveryQueues: routes.NewDeliveryQueuesRouter(container.DeliveryQueues),
	}
}

func (c ChiRouter) StartHttpServer(appStartTime time.Time) {
	port := viper.GetString("server.port")

	logging.Logger.Info(fmt.Sprintf("[STARTUP] Application started in %dms\n", time.Since(appStartTime).Milliseconds()))
	e := http.ListenAndServe(fmt.Sprintf(":%s", port), c.Router)
	if e != nil {
		panic(e)
	}
}

// FileServer is serving static files.
func FileServer(router chi.Router) {
	route := "/static"

	box := packr.NewBox("../../../resources/public")
	fs := http.FileServer(box)

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(route+"/", fs).ServeHTTP(w, r)
	})
}

func (c ChiRouter) Cors() func(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "issuer_id", "Content-Type", "Issuerid"},
		ExposedHeaders:     []string{"Content-Length", "Access-Control-Allow-Origin", "Content-Type", "issuer_id"},
		AllowCredentials:   true,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
		OptionsPassthrough: false,
	}).Handler
}

// func (c ChiRouter) Swagger() http.HandlerFunc {
// 	// Overwriting swagger generated data
// 	docs.SwaggerInfo.Title = global.Viper.GetString("app.name")
// 	docs.SwaggerInfo.Version = global.VersionInfo.Current
// 	docs.SwaggerInfo.Description = global.Viper.GetString("app.description")
// 	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", global.Viper.GetString("swagger.ip"), global.Viper.GetString("swagger.port"))

// 	// Creating a route to serve the generated JSON
// 	c.Router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		_, _ = w.Write(helpers.FetchSwagger(viper.GetBool("swagger.hideInternals") && (r.URL.Query().Get("internalsToken") != global.Viper.GetString("swagger.internalsToken"))))
// 	})

// 	swaggerUrl := "/swagger.json"
// 	if viper.GetBool("swagger.customJson.enabled") {
// 		swaggerUrl = viper.GetString("swagger.customJson.url")
// 	}
// 	// Informing swaggo the path to swagger and creating the handler
// 	return httpSwagger.Handler(httpSwagger.URL(swaggerUrl))
// }

func (c ChiRouter) CreateRoutes(serviceContainer domains.ServiceContainer) {
	routesContainer := NewRouterContainer(serviceContainer)
	c.Router.Use(c.Cors())
	c.Router.Mount("/debug", middleware2.Profiler())

	c.Router.Route("/v1", func(r chi.Router) {
		r.Mount("/metrics", promhttp.Handler())

		r.Group(func(logGroup chi.Router) {
			logGroup.Use(global.BillingInterceptors.HttpInterceptor)
			logGroup.Use(middleware.Logger)
			logGroup.Use(middleware.Recoverer)
			logGroup.Use(render.SetContentType(render.ContentTypeJSON))
			logGroup.Route("/internals", routesContainer.Internals.Route)
			logGroup.Route("/encryption", routesContainer.Encryption.Route)
			logGroup.Route("/webhooks", routesContainer.WebHooks.Route)
			logGroup.Route("/amqp", routesContainer.AMQP.Route)
			logGroup.Route("/events", routesContainer.Events.Route)
			logGroup.Route("/issuers", routesContainer.IssuerEvents.Route)
			logGroup.Route("/delivery-queues", routesContainer.DeliveryQueues.Route)
		})
		r.Group(func(managerGroup chi.Router) {
			managerGroup.Use(render.SetContentType(render.ContentTypeJSON))
			managerGroup.Route("/manager", routesContainer.Management.Route)
		})
	})

	// c.Router.Get("/swagger/*", c.Swagger())
	FileServer(c.Router)
}

func ConfigChi(servicesContainer domains.ServiceContainer, started time.Time) {
	chiRouter := ChiRouter{Router: chi.NewRouter()}
	chiRouter.CreateRoutes(servicesContainer)
	chiRouter.StartHttpServer(started)
}
