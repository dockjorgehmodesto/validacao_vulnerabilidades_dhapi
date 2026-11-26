package splunk

import (
	"database/sql"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"go.uber.org/zap/zapcore"
)

type DatabaseConnectionLog struct {
	AppLog
	Trace *DatabaseConnectionTrace `json:"trace"`
}

func (d DatabaseConnectionLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("@version", d.Version)
	enc.AddString("app_name", d.AppName)
	enc.AddString("app_profile", d.AppProfile)
	enc.AddString("type", d.Type)
	enc.AddString("level", d.Level) // TODO: remover depois da migração para o datadog
	// O Zap é inteligente: como o p.Trace já tem um MarshalLogObject,
	// nós adicionamos ele como um objeto aninhado usando AddObject!
	enc.AddObject("trace", d.Trace)
	return nil
}

func (obj *DatabaseConnectionLog) InitLog(stats sql.DBStats, mapKey string, uuid string) DatabaseConnectionLog {

	obj.Version = global.VersionInfo.Current
	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Level = LevelDebug
	obj.Type = TypeDatabaseConnection
	obj.Time = time.Now()
	obj.Trace = &DatabaseConnectionTrace{
		AppLogTrace: AppLogTrace{
			InsertedOnDate: time.Now(),
		},
		HealthcheckUUID: uuid,
		MapKey:          mapKey,
		IsApplicationDb: mapKey == "TENANT_PRINCIPAL",
		ConnectionsCount: ConnectionsCount{
			Idle:  stats.Idle,
			InUse: stats.InUse,
			Open:  stats.OpenConnections,
		},
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
		MaxOpenConnections: stats.MaxOpenConnections,
	}
	obj.Trace.Host = obj.Trace.GetHost()
	return *obj
}

func (obj DatabaseConnectionLog) Send(message string) {
	ShowLog(obj, obj.Level, obj.Type, message)
}

type DatabaseConnectionTrace struct {
	AppLogTrace
	HealthcheckUUID string `json:"healthcheckUuid"`

	MapKey           string           `json:"mapKey"`
	IsApplicationDb  bool             `json:"isApplicationDb"`
	ConnectionsCount ConnectionsCount `json:"connectionsCount"`

	WaitCount          int64         `json:"waitCount"`
	WaitDuration       time.Duration `json:"waitDuration"`
	MaxIdleClosed      int64         `json:"maxIdleClosed"`
	MaxLifetimeClosed  int64         `json:"maxLifetimeClosed"`
	MaxOpenConnections int           `json:"maxOpenConnections"`
}

func (d DatabaseConnectionTrace) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("healthcheckUuid", d.HealthcheckUUID)
	enc.AddString("mapKey", d.MapKey)
	enc.AddBool("isApplicationDb", d.IsApplicationDb)
	enc.AddObject("connectionsCount", d.ConnectionsCount)
	enc.AddInt64("waitCount", d.WaitCount)
	enc.AddDuration("waitDuration", d.WaitDuration)
	enc.AddInt64("maxIdleClosed", d.MaxIdleClosed)
	enc.AddInt64("maxLifetimeClosed", d.MaxLifetimeClosed)
	enc.AddInt("maxOpenConnections", d.MaxOpenConnections)
	return nil
}

type ConnectionsCount struct {
	Idle  int `json:"idle"`
	InUse int `json:"inUse"`
	Open  int `json:"open"`
}

func (c ConnectionsCount) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("idle", c.Idle)
	enc.AddInt("inUse", c.InUse)
	enc.AddInt("open", c.Open)
	return nil
}
