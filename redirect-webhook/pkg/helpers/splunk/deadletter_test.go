package splunk

import (
	"errors"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
)

func TestDeadLetterLog_InitLog(t *testing.T) {

	t.Run("assure that structure is initialized correctly", func(t *testing.T) {
		log := &DeadLetterLog{}
		log.InitLog("myTest", "mysub", 111)

		hostname, _ := os.Hostname()

		assert.Equal(t, NameApp, log.AppName)
		assert.Equal(t, Profile, log.AppProfile)
		assert.Equal(t, LevelInfo, log.Level)
		assert.NotNil(t, log.Trace)
		assert.NotEmpty(t, log.Trace)
		assert.Equal(t, hostname, log.Trace.Host)
		assert.Equal(t, 111, log.Trace.IssuerId)
		assert.Equal(t, "myTest", log.Trace.EventType)
		assert.Equal(t, "mysub", log.Trace.Subscription)
	})
}

func TestDeadLetterLog_SendCustomError(t *testing.T) {
	t.Run("sending error", func(t *testing.T) {
		SplunkEnabled = false

		log := &DeadLetterLog{}
		err := errors.New(random.String(15))
		errorDescription := random.String(15)
		errMetadata := MessagingErrorMetadata{
			IssuerId:     rand.Int(),
			EventType:    random.String(15),
			Subscription: random.String(15),
		}
		log.SendCustomError(errorDescription, err, errMetadata)

		hostname, _ := os.Hostname()
		assert.Equal(t, NameApp, log.AppName)
		assert.Equal(t, Profile, log.AppProfile)
		assert.Equal(t, LevelError, log.Level)
		assert.NotNil(t, log.Trace)
		assert.NotEmpty(t, log.Trace)
		assert.Equal(t, hostname, log.Trace.Host)
		assert.Equal(t, TypeDeadLetterQueueJob, log.Type)
		assert.Equal(t, errMetadata.IssuerId, log.Trace.IssuerId)
		assert.Equal(t, errMetadata.EventType, log.Trace.EventType)
		assert.Equal(t, errMetadata.Subscription, log.Trace.Subscription)
		assert.NotEmpty(t, log.Trace.InsertedOnDate)
		assert.Equal(t, errorDescription, log.Trace.Message)
	})
}

func TestDeadLetterLog_LogProcess(t *testing.T) {
	t.Run("logging process", func(t *testing.T) {
		SplunkEnabled = false

		log := &DeadLetterLog{}
		eventType := random.String(15)
		subscription := random.String(15)
		issuerId := rand.Int()
		log.InitLog(eventType, subscription, issuerId)

		log.LogProcess()
		assert.Equal(t, NameApp, log.AppName)
		assert.Equal(t, Profile, log.AppProfile)
		assert.Equal(t, LevelInfo, log.Level)
		assert.NotNil(t, log.Trace)
		assert.NotEmpty(t, log.Trace)
		assert.Equal(t, TypeDeadLetterQueueProcess, log.Type)
		assert.Equal(t, issuerId, log.Trace.IssuerId)
		assert.Equal(t, eventType, log.Trace.EventType)
		assert.Equal(t, subscription, log.Trace.Subscription)
		assert.NotEmpty(t, log.Trace.InsertedOnDate)
	})
}

func TestDeadLetterLog_LogJob(t *testing.T) {
	t.Run("logging job", func(t *testing.T) {
		SplunkEnabled = false

		log := &DeadLetterLog{}
		eventType := random.String(15)
		subscription := random.String(15)
		issuerId := rand.Int()
		log.InitLog(eventType, subscription, issuerId)

		log.LogJob()
		assert.Equal(t, NameApp, log.AppName)
		assert.Equal(t, Profile, log.AppProfile)
		assert.Equal(t, LevelInfo, log.Level)
		assert.NotNil(t, log.Trace)
		assert.NotEmpty(t, log.Trace)
		assert.Equal(t, TypeDeadLetterQueueJob, log.Type)
		assert.Equal(t, issuerId, log.Trace.IssuerId)
		assert.Equal(t, eventType, log.Trace.EventType)
		assert.Equal(t, subscription, log.Trace.Subscription)
		assert.NotEmpty(t, log.Trace.InsertedOnDate)
	})
}
