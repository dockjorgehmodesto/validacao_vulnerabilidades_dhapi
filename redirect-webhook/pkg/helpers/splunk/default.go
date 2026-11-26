package splunk

import (
	"bytes"
	"net"
	"os"
	"strings"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	NameApp string
	Profile string

	LogLevel          string
	SplunkHost        string
	LogSplunkLevel    string
	SplunkEnabled     bool
	LogSplunkToStdOut bool
	useTcp            bool
)

const (
	TypeHttp                   = "HTTP"
	TypeAMQP                   = "AMQP"
	TypeDeadLetterQueueJob     = "AMQP-DLQ-JOB"
	TypeDeadLetterQueueProcess = "AMQP-DLQ-PROCESS"
	TypeConnectionManagement   = "CONNECTIONMANAGEMENT"
	TypeDatabaseConnection     = "DATABASECONNECTION"
	TypeCache                  = "CACHE"
)

const (
	LevelError = "ERROR"
	LevelWarn  = "WARN"
	LevelInfo  = "INFO"
	LevelTrace = "TRACE"
	LevelDebug = "DEBUG"
)

func init() {
	SplunkEnabled = global.Viper.GetBool("log.splunk")
	NameApp = global.Viper.GetString("app.name")
	Profile = global.Viper.GetString("app.profile")
	SplunkHost = global.Viper.GetString("splunk.dial")
	LogSplunkToStdOut = global.Viper.GetBool("log.stdOut.enable")
	LogSplunkLevel = strings.ToUpper(global.Viper.GetString("log.level"))
	// splunk.InitSplunkLog(splunk.Config{
	// 	NameApp: NameApp,
	// 	Version: global.VersionInfo.Current,
	// 	Profile: Profile,
	// 	Host:    global.Viper.GetString("splunk.hec.host"),
	// 	Token:   global.Viper.GetString("splunk.hec.token"),
	// 	Tls:     global.Viper.GetBool("splunk.hec.skipTls"),
	// 	Source:  NameApp,
	// })
	useTcp = global.Viper.GetBool("splunk.useTcp")
	LogLevel = strings.ToUpper(global.Viper.GetString("log.level"))
	// go syncChannelToTcp()
}

type AppLog struct {
	Time       time.Time `json:"@time"`
	Version    string    `json:"@version"`
	Level      string    `json:"level"`
	AppName    string    `json:"app_name"`
	AppProfile string    `json:"app_profile"`
	Type       string    `json:"type"`
}

func (a AppLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("@version", a.Version)
	enc.AddString("app_name", a.AppName)
	enc.AddString("app_profile", a.AppProfile)
	enc.AddString("type", a.Type)
	enc.AddString("level", a.Level) // TODO: remover depois da migração para o datadog
	return nil
}

type AppLogTrace struct {
	IssuerId       int                    `json:"issuerId"`
	Host           string                 `json:"host"`
	InsertedOnDate time.Time              `json:"insertedOnDate"`
	Message        interface{}            `json:"message,omitempty"`
	Logs           interface{}            `json:"logs,omitempty"`
	Durations      map[string]interface{} `json:"durations,omitempty"`
}

func (a AppLogTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("issuerId", a.IssuerId)
	encoder.AddString("host", a.Host)
	// encoder.AddTime("insertedOnDate", a.InsertedOnDate)
	encoder.AddReflected("message", a.Message)
	if a.Logs != nil {
		encoder.AddReflected("logs", a.Logs)
	}
	if a.Durations != nil {
		encoder.AddReflected("durations", a.Durations)
	}
	return nil
}

func (s *AppLogTrace) GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname
}

type ConsoleLog struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	File    string    `json:"file"`
	Line    string    `json:"line"`
	Message string    `json:"message"`
}

func (c ConsoleLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddTime("time", c.Time)
	encoder.AddString("level", c.Level)
	encoder.AddString("file", c.File)
	encoder.AddString("line", c.Line)
	encoder.AddString("message", c.Message)
	return nil
}

type SplunkContext struct {
	Buf   bytes.Buffer
	Stack interface{}
}

var (
	splunkTcpConnection net.Conn
	// splunkTcpConnectionMutex = sync.Mutex{}
	splunkMessageChan chan []byte
)

// func SendToSplunk(value []byte, logLevel string, logType string) {
// 	if logType != TypeHttp {

// 		switch LogSplunkLevel {
// 		case LevelError:
// 			if logLevel != LevelError {
// 				return
// 			}
// 		case LevelWarn:
// 			if logLevel == LevelInfo || logLevel == LevelTrace {
// 				return
// 			}

// 		}
// 	}

// 	if LogSplunkToStdOut {
// 		fmt.Println(string(value))
// 	}

// 	go func() {
// 		if SplunkEnabled {
// 			if useTcp {
// 				// Uses socket sender
// 				splunkMessageChan <- value
// 			} else {
// 				// Uses Http Event Collector
// 				splunk.SendToSplunk(value)
// 			}
// 		}
// 	}()
// }

func ShowLog(value zapcore.ObjectMarshaler, logLevel string, logType string, message string) {
	if logType != TypeHttp {
		switch LogLevel {
		case LevelError:
			if logLevel != LevelError {
				return
			}
		case LevelWarn:
			if logLevel == LevelInfo || logLevel == LevelTrace {
				return
			}
		}
	}
	switch logLevel {
	case LevelInfo:
		logging.WrapperLogger.Info(message, zap.Inline(value))
	case LevelError:
		logging.WrapperLogger.Error(message, zap.Inline(value))
	case LevelWarn:
		logging.WrapperLogger.Warn(message, zap.Inline(value))
	default:
		logging.WrapperLogger.Debug(message, zap.Inline(value))
	}
}

// func syncChannelToTcp() {
// 	splunkMessageChan = make(chan []byte)
// 	for {
// 		var err error
// 		splunkMsg := <-splunkMessageChan
// 		if splunkTcpConnection == nil {
// 			splunkTcpConnection, err = net.Dial("tcp", SplunkHost)
// 			if err != nil {
// 				log.Errorf("Error connecting to splunk. Error: ", err.Error())
// 				logging.Logger.Error("[SPLUNK] could not connect", zap.Error(err))
// 				time.Sleep(15 * time.Second)
// 				continue
// 			} else {
// 				logging.Logger.Info("[SPLUNK] connecting")
// 			}
// 			hostname, _ := os.Hostname()
// 			splunkTcpConnection.Write(
// 				[]byte(
// 					fmt.Sprintf("{\"trace\": {\"message\": \"connecting to splunk\", \"host\": \"%s\"}, \"@time\": \"%s\", \"@version\": \"%s\", \"level\": \"%s\", \"app_name\": \"%s\", \"app_profile\": \"%s\", \"type\":\"%s\"}\n",
// 						hostname,
// 						time.Now().Format(time.RFC3339),
// 						global.VersionInfo.Current,
// 						LevelTrace,
// 						NameApp,
// 						Profile,
// 						"SPLUNKCONNECTION",
// 					),
// 				),
// 			)
// 		}

// 		_, err = splunkTcpConnection.Write(append(splunkMsg, []byte("\n")...))
// 		if err != nil {
// 			logging.Logger.Error("[SPLUNK] could not send message", zap.Error(err))
// 			err = splunkTcpConnection.Close()
// 			if err != nil {
// 				logging.Logger.Error("[SPLUNK] could not close connection", zap.Error(err))
// 			}
// 			splunkTcpConnection = nil
// 			// splunkMessageChan <- splunkMsg
// 		}
// 	}
// }
