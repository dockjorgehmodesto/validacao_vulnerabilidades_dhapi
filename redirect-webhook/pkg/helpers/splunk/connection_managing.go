package splunk

import (
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"go.uber.org/zap/zapcore"
)

type ConnectionManagementLog struct {
	AppLog
	Trace ConnectionManagementTrace `json:"trace"`
}

func (a *ConnectionManagementLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
	enc.AddString("@version", a.Version)
	enc.AddString("app_name", a.AppName)
	enc.AddString("app_profile", a.AppProfile)
	enc.AddString("type", a.Type)
	enc.AddString("level", a.Level) // TODO: remover depois da migração para o datadog
	// O Zap é inteligente: como o p.Trace já tem um MarshalLogObject,
	// nós adicionamos ele como um objeto aninhado usando AddObject!
	enc.AddObject("trace", a.Trace)
	return nil
}

type ConnectionManagementTrace struct {
	AppLogTrace
	Stacktrace   interface{} `json:"stacktrace"`
	ConnectionId string      `json:"connectionId"`
	EventType    string      `json:"eventType"`
}

func (c ConnectionManagementTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
	encoder.AddInt("issuerId", c.IssuerId)
	encoder.AddString("host", c.Host)
	encoder.AddTime("insertedonDate", c.InsertedOnDate)

	if c.Message != nil {
		encoder.AddReflected("message", c.Message)
	}
	if c.Logs != nil {
		encoder.AddReflected("logs", c.Logs)
	}
	if c.Durations != nil {
		encoder.AddReflected("durations", c.Durations)
	}
	if c.Stacktrace != nil {
		encoder.AddReflected("stacktrace", c.Stacktrace)
	}
	encoder.AddString("connectionId", c.ConnectionId)
	encoder.AddString("eventType", c.EventType)
	return nil
}

type MessageOpts struct {
	IssuerId     int    `json:"issuerId"`
	ConnectionId string `json:"connectionId"`
	EventType    string `json:"eventType"`
}

func (obj *ConnectionManagementLog) Setup() *ConnectionManagementLog {
	obj.Time = time.Now()
	obj.AppName = NameApp
	obj.Version = global.VersionInfo.Current
	obj.AppProfile = Profile
	obj.Type = TypeConnectionManagement
	obj.Trace = ConnectionManagementTrace{
		AppLogTrace: AppLogTrace{
			InsertedOnDate: time.Now(),
		},
	}
	obj.Trace.Host = obj.Trace.GetHost()
	return obj
}

func (obj ConnectionManagementLog) SendInfo(msg string) {
	obj.Setup()
	obj.Level = LevelInfo
	obj.Trace.Message = msg
	obj.send(msg)
}

func (obj ConnectionManagementLog) SendWarn(msg string, err error, opts *MessageOpts) {

	obj.Setup()
	obj.Level = LevelWarn
	obj.Trace.Message = msg
	obj.Trace.Stacktrace = err

	obj.updateOpts(opts)
	obj.send(msg)
}

func (obj *ConnectionManagementLog) updateOpts(opts *MessageOpts) {
	if opts != nil {
		obj.Trace.EventType = opts.EventType
		obj.Trace.ConnectionId = opts.ConnectionId
		obj.Trace.IssuerId = opts.IssuerId
	}
}

func (obj ConnectionManagementLog) SendError(msg string, err error, opts *MessageOpts) {

	obj.Setup()
	obj.Level = LevelError
	obj.Trace.Message = msg
	obj.Trace.Stacktrace = err
	obj.updateOpts(opts)
	obj.send(msg)
}

func (obj *ConnectionManagementLog) send(message string) {
	ShowLog(obj, obj.Level, obj.Type, message)
}
