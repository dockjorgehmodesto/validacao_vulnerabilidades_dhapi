package splunk

import (
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"github.com/eko/gocache/codec"
)

type CacheStatsLog struct {
	AppLog
	Trace *CacheTrace `json:"trace"`
}

func (obj *CacheStatsLog) WithStats(stats codec.Stats) CacheStatsLog {

	obj.AppName = NameApp
	obj.AppProfile = Profile
	obj.Level = LevelInfo
	obj.Type = TypeCache
	obj.Time = time.Now()
	obj.Trace = &CacheTrace{
		AppLogTrace: AppLogTrace{
			InsertedOnDate: time.Now(),
		},
		Stats: stats,
	}
	obj.Trace.Host = obj.Trace.GetHost()
	obj.Version = global.VersionInfo.Current
	return *obj
}

func (obj CacheStatsLog) Send() {
	// objJson, _ := json.Marshal(&obj)
	// SendToSplunk(objJson, obj.Level, obj.Type)

	ShowLog(obj, obj.Level, obj.Type, "cache stats")
}

type CacheTrace struct {
	AppLogTrace
	codec.Stats
}
