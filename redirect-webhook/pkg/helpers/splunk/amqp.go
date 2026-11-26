package splunk

import (
	"context"
	"fmt"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type AMQPLog struct {
	AppLog
	Trace      AmqpTrace `json:"trace"`
	Source     string    `json:"amqpSource"`
	BrokerType string    `json:"brokerType"`
}

func (a *AMQPLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
	enc.AddString("@version", a.Version)
	enc.AddString("app_name", a.AppName)
	enc.AddString("app_profile", a.AppProfile)
	enc.AddString("type", a.Type)
	enc.AddString("level", a.Level) // TODO: remover depois da migração para o datadog
	enc.AddString("amqpSource", a.Source)
	enc.AddString("brokerType", a.BrokerType)
	// O Zap é inteligente: como o p.Trace já tem um MarshalLogObject,
	// nós adicionamos ele como um objeto aninhado usando AddObject!
	enc.AddObject("trace", a.Trace)
	return nil
}

func NewAMQPLog() *AMQPLog {
	return &AMQPLog{}
}

type AmqpTrace struct {
	AppLogTrace

	Stacktrace      interface{}          `json:"stacktrace"`
	DurationSec     int64                `json:"durationSec"`
	WebHookResponse helpers.HttpResponse `json:"webHookResponse"`
	MessageId       string               `json:"requestId"`
	EventType       string               `json:"eventType"`
	EncryptionKeyId string               `json:"encryptionKeyId"`
	WebHookId       string               `json:"webHookId"`
	Subscription    string               `json:"subscription,omitEmpty"`
	DeliveryCount   int                  `json:"deliveryCount"`
	ConnectionId    string               `json:"connectionId"`
}

func (a AmqpTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("issuerId", a.IssuerId)
	encoder.AddString("host", a.Host)
	encoder.AddTime("insertedonDate", a.InsertedOnDate)
	encoder.AddReflected("message", a.Message)
	if a.Logs != nil {
		encoder.AddReflected("logs", a.Logs)
	}
	if a.Durations != nil {
		encoder.AddReflected("durations", a.Durations)
	}
	encoder.AddReflected("stacktrace", a.Stacktrace)
	encoder.AddInt64("durationsSec", a.DurationSec)
	encoder.AddReflected("webhookResponse", a.WebHookResponse)
	encoder.AddString("requestId", a.MessageId)
	encoder.AddString("eventType", a.EventType)
	encoder.AddString("encryptionKeyId", a.EncryptionKeyId)
	encoder.AddString("webHookId", a.WebHookId)
	encoder.AddString("subscription", a.Subscription)
	encoder.AddInt("deliveryCount", a.DeliveryCount)
	encoder.AddString("connectionId", a.ConnectionId)
	return nil
}

func (obj *AMQPLog) InitServiceBusLog(msgId string, eventType string, subscription string, issuerId int, deliveryCount int, connectionId string) {

	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Type = TypeAMQP
	obj.BrokerType = "SERVICEBUS"
	obj.Level = LevelInfo
	obj.Version = global.VersionInfo.Current
	obj.Trace = AmqpTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       issuerId,
			InsertedOnDate: time.Now(),
		},
		MessageId:     msgId,
		EventType:     eventType,
		Subscription:  subscription,
		DeliveryCount: deliveryCount,
		ConnectionId:  connectionId,
	}
	obj.Trace.Host = obj.Trace.GetHost()
}

func (obj *AMQPLog) InitRabbitMQLog(msgId string, eventType string, issuerId int, deliveryCount int, connectionId string) {

	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Type = TypeAMQP
	obj.BrokerType = "RABBITMQ"
	obj.Level = LevelInfo
	obj.Version = global.VersionInfo.Current
	obj.Trace = AmqpTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       issuerId,
			InsertedOnDate: time.Now(),
		},
		MessageId:     msgId,
		EventType:     eventType,
		DeliveryCount: deliveryCount,
		ConnectionId:  connectionId,
	}
	obj.Trace.Host = obj.Trace.GetHost()
}

type MessagingErrorMetadata struct {
	IssuerId     int
	EventType    string
	Subscription string
	ConnectionId string
	AmqpSource   string
}

func (obj *AMQPLog) SendCustomError(msg string, err error, metadata MessagingErrorMetadata) {
	obj.AppName = NameApp
	obj.Version = global.VersionInfo.Current
	obj.AppProfile = Profile
	obj.Type = TypeAMQP
	trace := AmqpTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       metadata.IssuerId,
			InsertedOnDate: time.Now(),
		},
		EventType:    metadata.EventType,
		Subscription: metadata.Subscription,
		ConnectionId: metadata.ConnectionId,
	}

	obj.Time = time.Now()
	obj.Level = LevelError
	trace.Message = msg
	trace.Stacktrace = fmt.Sprintf("%+v", err)
	obj.Trace = trace
	obj.Trace.Host = obj.Trace.GetHost()
	// objJson, _ := json.Marshal(&obj)
	// SendToSplunk(objJson, obj.Level, obj.Type)
	ShowLog(obj, obj.Level, obj.Type, msg)
}

func (obj *AMQPLog) send(ctx context.Context, response *helpers.HttpResponse, keyId string, err error, level string, webhookId string, message string) {
	obj.Trace.DurationSec = time.Since(obj.Trace.InsertedOnDate).Milliseconds()
	obj.Time = time.Now()
	if response != nil {
		obj.Trace.WebHookResponse = *response
	}
	obj.Level = level
	obj.fetchMetadataFromContext(ctx)
	if err != nil {
		obj.Trace.Message = err.Error()
		obj.Trace.Stacktrace = errors.Cause(err).Error()
	}

	obj.Trace.EncryptionKeyId = keyId
	obj.Trace.WebHookId = webhookId
	obj.Trace.Host = obj.Trace.GetHost()

	ShowLog(obj, obj.Level, obj.Type, message)
}

func (obj *AMQPLog) SendCustomMessage(level string, msg string) {
	obj.Trace.Host = obj.Trace.GetHost()
	obj.Time = time.Now()
	obj.Level = level
	obj.Trace.Message = msg
	// objJson, _ := json.Marshal(&obj)
	// SendToSplunk(objJson, obj.Level, obj.Type)
	ShowLog(obj, obj.Level, obj.Type, msg)
}

func (obj *AMQPLog) SendWarnForWebhook(ctx context.Context, response *helpers.HttpResponse, keyId string, err error) {
	obj.send(ctx, response, keyId, err, LevelWarn, "", "SendWarnForWebhook")
}

func (obj *AMQPLog) SendErrorForWebhook(ctx context.Context, response *helpers.HttpResponse, keyId string, err error) {
	obj.send(ctx, response, keyId, err, LevelError, "", "SendErrorForWebhook")
}

func (obj *AMQPLog) SendInfoForWebhook(ctx context.Context, response *helpers.HttpResponse, keyId string, webHookId string) {
	obj.send(ctx, response, keyId, nil, LevelInfo, webHookId, "SendInfoForWebhook")
}

func (obj *AMQPLog) fetchMetadataFromContext(ctx context.Context) {
	if ctx != nil {
		splunkContextPointer := ctx.Value("splunkContext" + obj.Trace.MessageId)
		if splunkContextPointer != nil {
			splunkContext := splunkContextPointer.(*SplunkContext)
			var durations map[string]interface{}
			if splunkContext.Buf.Len() > 0 {
				var logSplunk []ConsoleLog
				logSplunk, durations = FormatLog(splunkContext.Buf.String())
				obj.Trace.Logs = logSplunk
				obj.Trace.Durations = durations
			}
			obj.Trace.Message = splunkContext.Stack
		}
	}
}
