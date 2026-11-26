package global

import (
	"bytes"
	"fmt"
	"strings"

	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/gobuffalo/packr"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {

	logging.Logger.Info("[VIPER] Loading context")
	configName := "application"

	viperSetup := viper.GetViper()
	viperSetup.SetConfigType("yml")
	viperSetup.SetConfigName(configName)
	viperSetup.AllowEmptyEnv(true)
	box := packr.NewBox("../../../resources")
	s, err := box.Find(configName + ".yml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	viperSetup.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viperSetup.AutomaticEnv()
	if err := viper.ReadConfig(bytes.NewBuffer(s)); err != nil {
		logging.Logger.Fatal("[VIPER] Error reading configuration files", zap.Error(err))
	}
	Viper = viperSetup

	ShowSQL = Viper.GetBool("database.showSQL")
}
