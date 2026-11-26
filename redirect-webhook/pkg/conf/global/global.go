package global

import (
	"conductor.com.br/pier/hooks/pkg"
	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	"github.com/spf13/viper"
)

var (
	//Viper godoc
	Viper *viper.Viper

	ShowSQL = false

	VersionInfo pkg.VersionData

	BillingInterceptors billing.Interceptors

	BillingLogger billing.Logger
)
