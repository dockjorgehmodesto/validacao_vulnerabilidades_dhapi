package splunk

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/gofrs/uuid"
	"go.uber.org/zap/zapcore"
)

type HttpLog struct {
	AppLog
	Verb      string      `json:"verb"`
	Url       string      `json:"url"`
	RequestID interface{} `json:"requestID"`
	Trace     HttpTrace   `json:"trace"`
}

func (h *HttpLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	// Extraindo os dados do AppLog embutido
	encoder.AddString("@version", h.Version)
	encoder.AddString("app_name", h.AppName)
	encoder.AddString("app_profile", h.AppProfile)
	encoder.AddString("type", h.Type)
	encoder.AddString("level", h.Level) // TODO: remover depois da migração para o datadog
	if h.RequestID != nil {
		encoder.AddReflected("requestID", h.RequestID)
	}
	encoder.AddString("verb", h.Verb)
	encoder.AddString("url", h.Url)
	encoder.AddObject("trace", h.Trace)
	return nil
}

type HttpTrace struct {
	AppLogTrace
	DurationSec         int64       `json:"durationSec"`
	ReceivedFromAddress string      `json:"receivedFromAddress"`
	ResultStatus        int         `json:"resultStatus"`
	Pattern             interface{} `json:"pattern"`
}

func (p HttpTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("issuerId", p.IssuerId)
	encoder.AddString("host", p.Host)
	if p.Message != nil {
		encoder.AddReflected("message", p.Message)
	}
	if p.Logs != nil {
		encoder.AddReflected("logs", p.Logs)
	}

	if p.Durations != nil {
		encoder.AddReflected("durations", p.Durations)
	}
	encoder.AddInt64("durationSec", p.DurationSec)
	encoder.AddString("receivedFromAddress", p.ReceivedFromAddress)
	encoder.AddInt("resultStatus", p.ResultStatus)
	if p.Pattern != nil {
		encoder.AddReflected("pattern", p.Pattern)
	}
	return nil
}

func (obj *HttpLog) InitRequestLog(r *http.Request) {
	issuer, _ := strconv.Atoi(r.Header.Get("issuer_id"))
	requestId := r.Header.Get("request_id")
	if requestId == "" {
		Uuid, _ := uuid.NewV4()
		requestId = Uuid.String()
	}

	obj.Version = global.VersionInfo.Current
	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Type = TypeHttp
	obj.Verb = r.Method
	obj.RequestID = requestId
	obj.Url = r.Host + r.RequestURI
	obj.Trace = HttpTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       issuer,
			InsertedOnDate: time.Now(),
		},
		ReceivedFromAddress: r.RemoteAddr,
		Pattern:             getRoutePattern(r),
	}
	obj.Trace.Host = obj.Trace.GetHost()
}

func (obj *HttpLog) SetLevelByStatusCode(statusCode int) {
	obj.Level = GetHttpLevelByStatus(statusCode)
}

func GetHttpLevelByStatus(statusCode int) string {
	if statusCode < 300 {
		return LevelInfo
	} else if statusCode >= 300 && statusCode < 500 {
		return LevelWarn
	} else {
		return LevelError
	}
}

func (obj *HttpLog) FetchMetadataFromContext(ctx context.Context) {
	splunkContextPointer := ctx.Value("splunkContext")
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

func (obj *HttpLog) Send(message string) {
	ShowLog(obj, obj.Level, obj.Type, message)
}

func getRoutePattern(r *http.Request) string {
	result := ""
	stringSlice := strings.Split(r.RequestURI, "/")
	for _, v := range stringSlice {
		if v != "" {
			if valid, err := helpers.ValidatorUUID(v); err != nil || !valid {
				result += "/" + v
			} else {
				result += "/{id}"
			}
		}
	}

	return result
}
