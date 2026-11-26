package conf

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupBilling() {
	applicationData, logOptions := getDataWithEnvironment()

	interceptor, err := billing.NewInterceptor(applicationData, logOptions)

	if err != nil {
		logging.Logger.Fatal("[BILLING] Could not create billing interceptor", zap.Error(err))
	}

	global.BillingInterceptors = interceptor

	logger, err := billing.NewLogger(applicationData, logOptions)

	if err != nil {
		logging.Logger.Fatal("[BILLING] Could not create billing logger", zap.Error(err))
	}

	global.BillingLogger = logger
}

func getDataWithEnvironment() (billing.ApplicationData, billing.LogOptions) {
	applicationData := billing.ApplicationData{
		Name:    global.Viper.GetString("app.name"),
		Version: global.Viper.GetString("app.version"),
	}
	logOptions := billing.LogOptions{
		Enabled: global.Viper.GetBool("log.billing.enabled"),
		TcpOptions: billing.TcpOptions{
			Enabled:          global.Viper.GetBool("log.billing.tcp.enabled"),
			ConnectionString: global.Viper.GetString("log.billing.tcp.connectionstring"),
		},
		FileOptions: billing.FileOptions{
			Enabled:    global.Viper.GetBool("log.billing.file.enabled"),
			Path:       global.Viper.GetString("log.billing.file.path"),
			MaxSizeMb:  global.Viper.GetInt("log.billing.file.maxsizemb"),
			MaxBackups: global.Viper.GetInt("log.billing.file.maxbackups"),
			MaxAgeDays: global.Viper.GetInt("log.billing.file.maxagedays"),
		},
		LogLevel: zapcore.InfoLevel,
	}
	return applicationData, logOptions
}
