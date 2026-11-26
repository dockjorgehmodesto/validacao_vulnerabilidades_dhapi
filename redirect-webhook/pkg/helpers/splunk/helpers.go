package splunk

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"strings"
	"time"
)

func LogElapsedTime(key string, initTime time.Time) log.JSON {
	result := log.JSON{}
	result = make(map[string]interface{})
	result[key] = time.Since(initTime).Milliseconds()
	return result
}

func FormatLog(buf string) ([]ConsoleLog, map[string]interface{}) {
	stackLog := []ConsoleLog{}
	nova := strings.SplitAfter(buf, "\n")
	jsonMap := make(map[string]interface{})
	durations := make(map[string]interface{})

	for _, v := range nova[:len(nova)-1] {
		objLog := ConsoleLog{}
		err := json.Unmarshal([]byte(v), &objLog)
		if err != nil {
			objLog.Level = "INFO"
			objLog.Message = strings.Trim(v, "\n")
		}

		// Parsing durations
		_ = json.Unmarshal([]byte(v), &jsonMap)
		for k, v := range jsonMap {
			if strings.Index(k, "elapsed") == 0 {
				durations[k] = v
			}
		}

		// Only add to result slice if message isn't empty
		if objLog.Message != "" {
			stackLog = append(stackLog, objLog)
		}
	}
	return stackLog, durations
}
