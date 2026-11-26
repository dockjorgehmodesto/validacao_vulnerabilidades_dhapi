package conf

import (
	"conductor.com.br/pier/hooks/pkg"
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/gobuffalo/packr"
	"go.uber.org/zap"
)

func ConfigVersion() {

	box := packr.NewBox("../../resources")
	bytes, e := box.Find(".version")

	if e != nil {
		logging.Logger.Fatal("[VERSION] could not read version file", zap.Error(e))
	}

	versionFile, e := pkg.ParseFile(bytes)
	if e != nil {
		logging.Logger.Fatal("[VERSION] could not parse version file", zap.Error(e))
	}
	global.VersionInfo = versionFile.Version
}
