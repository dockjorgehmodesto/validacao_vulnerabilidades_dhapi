package splunk

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pkg/errors"
)

type DeliveryLog struct {
	AppLog
	Trace DeliveryTrace `json:"trace"`
}

func (p *DeliveryLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
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

type DeliveryTrace struct {
	AppLogTrace

	Stacktrace        interface{} `json:"stacktrace,omitempty"`
	Duration          int64       `json:"duration"`
	MessageExternalId string      `json:"requestId"`
	MessageInternalId string      `json:"internalId,omitempty"`
	GroupId           string      `json:"groupId,omitempty"`

	StatusCode      int                     `json:"statusCode,omitempty"`
	Response        string                  `json:"response,omitempty"`
	RequestDuration int64                   `json:"requestDuration,omitempty"`
	Metadata        *DeliveryMetadata       `json:"metadata,omitempty"`
	Processor       *DeliveryProcessorTrace `json:"processor,omitempty"`
	BytesThroughput int64                   `json:"bytesThroughput,omitempty"`
}

func (d DeliveryTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("issuerId", d.IssuerId)
	encoder.AddString("host", d.Host)
	if d.Message != nil {
		encoder.AddReflected("message", d.Message)
	}

	if d.Logs != nil {
		encoder.AddReflected("logs", d.Logs)

	}
	if d.Durations != nil {
		encoder.AddReflected("durations", d.Durations)
	}
	if d.Stacktrace != nil {
		encoder.AddReflected("stacktrace", d.Stacktrace)
	}
	encoder.AddInt64("duration", d.Duration)
	encoder.AddString("requestId", d.MessageExternalId)
	encoder.AddString("groupId", d.GroupId)
	encoder.AddInt("statusCode", d.StatusCode)
	if d.Response != "" {
		encoder.AddString("response", d.Response)
	}
	encoder.AddInt64("requestDuration", d.RequestDuration)
	encoder.AddObject("metadata", d.Metadata)
	return nil
}

type DeliveryMetadata struct {
	EventType          string     `json:"eventType"`
	WebHookId          string     `json:"webHookId"`
	DeliveryCount      int        `json:"deliveryCount"`
	ConnectionId       string     `json:"connectionId"`
	KeyVaultId         string     `json:"keyVaultId"`
	ProcessedAt        time.Time  `json:"processedAt"`
	ProcessDelayMillis int64      `json:"processDelayMillis"`
	Host               string     `json:"host,omitempty"`
	GeneratedAt        *time.Time `json:"generatedAt,omitempty"`
	OriginDelayMillis  *int64     `json:"originDelayMillis,omitempty"`
	PrimaryKey         uint64     `json:"primaryKey,omitempty"`
	Operation          string     `json:"operation,omitempty"`
	PublishAt          *time.Time `json:"publishAt,omitempty"`
	PublishDelayMillis *int64     `json:"publishDelayMillis,omitempty"`
}

func (d DeliveryMetadata) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("deliveryCount", d.DeliveryCount)
	encoder.AddInt64("queued", d.ProcessDelayMillis)
	encoder.AddString("eventType", d.EventType)
	encoder.AddString("webHookId", d.WebHookId)
	encoder.AddString("keyVaultId", d.KeyVaultId)
	encoder.AddTime("processedAt", d.ProcessedAt)
	encoder.AddString("host", d.Host)
	if d.GeneratedAt != nil {
		encoder.AddTime("generatedAt", *d.GeneratedAt)
	}
	encoder.AddUint64("primaryKey", d.PrimaryKey)
	encoder.AddString("operation", d.Operation)
	if d.PublishAt != nil {
		encoder.AddTime("publishAt", *d.PublishAt)
	}
	return nil
}

type DeliveryProcessorTrace struct {
	InsertedOnDate time.Time              `json:"insertedOnDate"`
	DurationSec    int64                  `json:"durationSec"`
	Durations      map[string]interface{} `json:"durations,omitempty"`
}

func NewDeliveryLog(issuerId int) *DeliveryLog {
	hostname, _ := os.Hostname()
	return &DeliveryLog{
		AppLog: AppLog{
			Time:       time.Now(),
			Version:    global.VersionInfo.Current,
			Level:      LevelInfo,
			AppName:    NameApp,
			AppProfile: Profile,
			Type:       "WEBHOOKDELIVERY",
		},
		Trace: DeliveryTrace{
			AppLogTrace: AppLogTrace{
				IssuerId:       issuerId,
				Host:           hostname,
				InsertedOnDate: time.Now(),
			},
		},
	}
}

func (d *DeliveryLog) WithTrace(trace DeliveryTrace) *DeliveryLog {
	d.Trace = DeliveryTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       d.Trace.IssuerId,
			Host:           d.Trace.Host,
			InsertedOnDate: d.Trace.InsertedOnDate,
			Message:        trace.Message,
			Logs:           trace.Logs,
			Durations:      trace.Durations,
		},
		Stacktrace:        trace.Stacktrace,
		Duration:          trace.Duration,
		MessageExternalId: trace.MessageExternalId,
		MessageInternalId: trace.MessageInternalId,
		GroupId:           trace.GroupId,
		Metadata:          trace.Metadata,
		Processor:         trace.Processor,
	}

	if trace.Metadata != nil && trace.Metadata.GeneratedAt != nil {
		milliseconds := time.Since(*trace.Metadata.GeneratedAt).Milliseconds()
		d.Trace.Metadata.OriginDelayMillis = &milliseconds
	}

	if trace.Metadata != nil && trace.Metadata.PublishAt != nil {
		pubMilliseconds := time.Since(*trace.Metadata.PublishAt).Milliseconds()
		d.Trace.Metadata.PublishDelayMillis = &pubMilliseconds
	}
	return d
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (d *DeliveryLog) WithError(err error) *DeliveryLog {
	if err != nil {
		if st, ok := err.(stackTracer); ok {
			d.Trace.Stacktrace = fmt.Sprintf("%+v", st.StackTrace()[0:2])
		} else if urlErr, ok := err.(*url.Error); ok {
			d.Trace.Stacktrace = struct {
				Timeout   bool   `json:"timeout"`
				Temporary bool   `json:"temporary"`
				Url       string `json:"url"`
				Cause     string `json:"cause"`
				Op        string `json:"op"`
			}{
				Timeout:   urlErr.Timeout(),
				Temporary: urlErr.Temporary(),
				Url:       urlErr.URL,
				Cause:     urlErr.Unwrap().Error(),
				Op:        urlErr.Op,
			}
		} else {
			d.Trace.Stacktrace = err
		}

	}
	return d
}

func (d *DeliveryLog) WithLevel(level string) *DeliveryLog {
	d.Level = level
	return d
}

func (d *DeliveryLog) WithResponse(response *helpers.HttpResponse) *DeliveryLog {
	// d.Trace.Response = response.Body
	d.Trace.BytesThroughput = response.TotalBytesThroughput
	d.Trace.StatusCode = response.Status
	d.Trace.RequestDuration = response.Duration
	return d
}

func (d *DeliveryLog) WithMessage(message string) *DeliveryLog {
	d.Trace.Message = message
	return d
}

func (d *DeliveryLog) Send(ctx context.Context, message string) {
	d.Trace.Duration = time.Since(d.Trace.InsertedOnDate).Milliseconds()
	d.Time = time.Now()
	d.fetchMetadataFromContext(ctx)

	// objJson, _ := json.Marshal(&d)
	logging.Logger.Debug("[DELIVERER]", zap.Object("event", d.Trace))
	ShowLog(d, d.AppLog.Level, d.AppLog.Type, message)
}

func (d *DeliveryLog) fetchMetadataFromContext(ctx context.Context) {
	if ctx != nil {
		splunkContextPointer := ctx.Value("deliveryContext" + d.Trace.MessageExternalId)
		if splunkContextPointer != nil {
			splunkContext := splunkContextPointer.(*SplunkContext)
			var durations map[string]interface{}
			if splunkContext.Buf.Len() > 0 {
				var logSplunk []ConsoleLog
				logSplunk, durations = FormatLog(splunkContext.Buf.String())
				d.Trace.Logs = logSplunk
				d.Trace.Durations = durations
			}
		}
	}
}
