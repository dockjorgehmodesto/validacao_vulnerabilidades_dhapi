package splunk

import (
	"context"
	"os"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/timer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DurationsKeyWebhook  = "elapsedFetchingWebhook"
	DurationsKeyKeyvault = "elapsedFetchingKey"
	DurationsKeyRabbitMQ = "elapsedFetchingRabbitMQ"
)

type processorsLog struct {
	AppLog
	Trace   ProcessorTrace `json:"trace"`
	Shipped bool           `json:"shipped"`
}

func (p *processorsLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
	// enc.AddTime("@time", p.Time)
	enc.AddString("@version", p.Version)
	enc.AddString("app_name", p.AppName)
	enc.AddString("app_profile", p.AppProfile)
	enc.AddString("type", p.Type)
	enc.AddString("level", p.Level) // TODO: remover depois da migração para o datadog
	// O Zap é inteligente: como o p.Trace já tem um MarshalLogObject,
	// nós adicionamos ele como um objeto aninhado usando AddObject!
	enc.AddObject("trace", p.Trace)
	return nil
}

type ProcessorTrace struct {
	AppLogTrace

	Stacktrace         interface{} `json:"stacktrace,omitempty"`
	DurationSec        int64       `json:"durationSec"`
	MessageExternalId  string      `json:"requestId"`
	MessageInternalId  string      `json:"internalId,omitempty"`
	EventType          string      `json:"eventType"`
	EncryptionKeyId    string      `json:"encryptionKeyId"`
	WebHookId          string      `json:"webHookId"`
	DeliveryCount      int         `json:"deliveryCount"`
	ConnectionId       string      `json:"connectionId"`
	BrokerType         string      `json:"brokerType"`
	GroupId            string      `json:"groupId"`
	GeneratedAt        *time.Time  `json:"generatedAt,omitempty"`
	OriginDelayMillis  *int64      `json:"originDelayMillis,omitempty"`
	PrimaryKey         uint64      `json:"primaryKey,omitempty"`
	Operation          string      `json:"operation,omitempty"`
	PublishAt          *time.Time  `json:"publishAt,omitempty"`
	PublishDelayMillis *int64      `json:"publishDelayMillis,omitempty"`
}

func (p ProcessorTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("issuerId", p.IssuerId)
	encoder.AddString("host", p.Host)
	if p.Message != nil {
		encoder.AddReflected("message", p.Message)
	}

	if p.Logs != nil {
		encoder.AddReflected("logs", p.Logs)

	}
	if p.Stacktrace != nil {
		encoder.AddReflected("stacktrace", p.Stacktrace)
	}
	encoder.AddInt("deliveryCount", p.DeliveryCount)
	encoder.AddInt64("durationSec", p.DurationSec)
	encoder.AddString("requestId", p.MessageExternalId)
	encoder.AddString("groupId", p.GroupId)
	encoder.AddString("eventType", p.EventType)
	encoder.AddString("webHookId", p.WebHookId)
	encoder.AddString("keyVaultId", p.EncryptionKeyId)
	encoder.AddString("brokerType", p.BrokerType)
	encoder.AddUint64("primaryKey", p.PrimaryKey)
	encoder.AddString("operation", p.Operation)
	if p.GeneratedAt != nil {
		encoder.AddTime("generatedAt", *p.GeneratedAt)
	}
	if p.PublishAt != nil {
		encoder.AddTime("publishAt", *p.PublishAt)
	}
	return nil
}

func NewProcessorsLog(issuerId int) *processorsLog {
	hostname, _ := os.Hostname()
	return &processorsLog{
		AppLog: AppLog{
			Time:       time.Now(),
			Version:    global.VersionInfo.Current,
			Level:      LevelInfo,
			AppName:    NameApp,
			AppProfile: Profile,
			Type:       "MESSAGEPROCESSOR",
		},
		Trace: ProcessorTrace{
			AppLogTrace: AppLogTrace{
				IssuerId:       issuerId,
				Host:           hostname,
				InsertedOnDate: time.Now(),
			},
		},
	}
}

func (p *processorsLog) WithTrace(trace ProcessorTrace) *processorsLog {
	p.Trace = ProcessorTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       p.Trace.IssuerId,
			Host:           p.Trace.Host,
			InsertedOnDate: p.Trace.InsertedOnDate,
			Message:        trace.Message,
			Logs:           trace.Logs,
			Durations:      trace.Durations,
		},
		Stacktrace:        trace.Stacktrace,
		DurationSec:       trace.DurationSec,
		MessageExternalId: trace.MessageExternalId,
		MessageInternalId: trace.MessageInternalId,
		EventType:         trace.EventType,
		EncryptionKeyId:   trace.EncryptionKeyId,
		WebHookId:         trace.WebHookId,
		DeliveryCount:     trace.DeliveryCount,
		ConnectionId:      trace.ConnectionId,
		BrokerType:        trace.BrokerType,
		GroupId:           trace.GroupId,
		GeneratedAt:       trace.GeneratedAt,
		PrimaryKey:        trace.PrimaryKey,
		Operation:         trace.Operation,
		PublishAt:         trace.PublishAt,
	}

	if trace.GeneratedAt != nil {
		milliseconds := time.Since(*trace.GeneratedAt).Milliseconds()
		p.Trace.OriginDelayMillis = &milliseconds
	}

	if trace.PublishAt != nil {
		pubMilliseconds := time.Since(*trace.PublishAt).Milliseconds()
		p.Trace.PublishDelayMillis = &pubMilliseconds
	}
	return p
}

func (p *processorsLog) WithError(err error, level string) *processorsLog {
	if err != nil {
		p.Trace.Message = err.Error()
		p.Trace.Stacktrace = errors.Cause(err).Error()
	}
	p.Level = level
	return p
}

func (p *processorsLog) WithLevel(level string) *processorsLog {
	p.Level = level
	return p
}

func (p *processorsLog) Send(ctx context.Context, message string) {
	p.Trace.DurationSec = time.Since(p.Trace.InsertedOnDate).Milliseconds()
	p.Time = time.Now()
	p.Trace.Durations = timer.GetTracker(ctx)

	logging.Logger.Debug("[PROCESSOR]", zap.Object("event", p.Trace))

	ShowLog(p, p.Level, p.Type, message)
}
