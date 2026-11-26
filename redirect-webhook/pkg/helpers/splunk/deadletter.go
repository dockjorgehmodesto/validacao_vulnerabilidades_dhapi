package splunk

import (
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"github.com/pkg/errors"
)

type DeadLetterLog struct {
	AppLog
	Trace *DeadLetterTrace `json:"trace"`
}

// Structure with information related
type DeadLetterTrace struct {
	AppLogTrace

	Stacktrace   interface{} `json:"stacktrace"`
	DurationSec  int64       `json:"durationSec"`
	EventType    string      `json:"eventType"`
	Subscription string      `json:"subscription"`
}

// Initializes the log structure for the dead letter context
func (obj *DeadLetterLog) InitLog(eventType string, subscription string, issuerId int) {

	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Level = LevelInfo
	obj.Trace = &DeadLetterTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       issuerId,
			InsertedOnDate: time.Now(),
		},
		EventType:    eventType,
		Subscription: subscription,
	}
	obj.Trace.Host = obj.Trace.GetHost()
}

// Logs a customized error based on the necessity
func (obj *DeadLetterLog) SendCustomError(msg string, err error, metadata MessagingErrorMetadata) {
	wrappedErr := errors.Wrap(err, msg)
	obj.AppName = NameApp
	obj.Version = global.VersionInfo.Current
	obj.AppProfile = Profile
	obj.Type = TypeDeadLetterQueueJob
	trace := DeadLetterTrace{
		AppLogTrace: AppLogTrace{
			IssuerId:       metadata.IssuerId,
			InsertedOnDate: time.Now(),
		},
		EventType:    metadata.EventType,
		Subscription: metadata.Subscription,
	}
	trace.Host = trace.GetHost()

	obj.Time = time.Now()
	obj.Level = LevelError
	trace.Message = msg
	trace.Stacktrace = errors.Cause(wrappedErr).Error()
	obj.Trace = &trace

	ShowLog(obj, obj.Level, obj.Type, msg)
}

// Logs a single issue or message inside a dead letter re-send process
func (obj *DeadLetterLog) LogJob() {
	obj.Trace.DurationSec = time.Since(obj.Trace.InsertedOnDate).Milliseconds()
	obj.Time = time.Now()
	obj.Level = LevelInfo
	obj.Type = TypeDeadLetterQueueJob

	ShowLog(obj, obj.Level, obj.Type, "log job")
}

// Indicates a Dead Letter Re-send process is done, and how did it perform
func (obj *DeadLetterLog) LogProcess() {
	obj.Trace.DurationSec = time.Since(obj.Trace.InsertedOnDate).Milliseconds()
	obj.Time = time.Now()
	obj.Level = LevelInfo
	obj.Type = TypeDeadLetterQueueProcess

	ShowLog(obj, obj.Level, obj.Type, "log process")
}
