package splunk

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"conductor.com.br/pier/hooks/pkg/conf/global"
// 	billing "dev.azure.com/conductortech/billing/_git/middleware-golang.git"
// )

// func init() {
// 	go func() {
// 		for err := range billing.GetErrChan() {
// 			if err != nil {
// 				hostname, _ := os.Hostname()
// 				logData := fmt.Sprintf("{\"trace\": {\"stacktrace\": \"%+v\", \"message\": \"connecting to logstash\", \"host\": \"%s\"}, \"@time\": \"%s\", \"@version\": \"%s\", \"level\": \"%s\", \"app_name\": \"%s\", \"app_profile\": \"%s\", \"type\":\"%s\"}\n",
// 					err,
// 					hostname,
// 					time.Now().Format(time.RFC3339),
// 					global.VersionInfo.Current,
// 					LevelTrace,
// 					NameApp,
// 					Profile,
// 					"LOGSTASHCONNECTION",
// 				)
// 				SendToSplunk([]byte(logData), LevelTrace, "LOGSTASHCONNECTION")
// 			}
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }
