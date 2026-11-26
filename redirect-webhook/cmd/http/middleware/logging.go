package middleware

import (
	bytes2 "bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"conductor.com.br/pier/hooks/pkg/helpers/splunk"
	"github.com/labstack/gommon/log"
)

type RequestLog struct {
	Status   *int
	Response *string
	Request  *string
}

func (r *RequestLog) Print() {
	bytes, _ := json.Marshal(r)
	fmt.Println(string(bytes))
}

type ResponseWrapper struct {
	http.ResponseWriter
	RequestLog *RequestLog
	Context    context.Context
}

func (r ResponseWrapper) Write(b []byte) (int, error) {
	value := string(b)
	r.RequestLog.Response = &value
	return r.ResponseWriter.Write(b)
}

func (r ResponseWrapper) WriteHeader(statusCode int) {
	r.RequestLog.Status = &statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseWrapper) ExecuteAfterResponding(msg *splunk.HttpLog) {
	msg.SetLevelByStatusCode(*r.RequestLog.Status)

	msg.Trace.DurationSec = time.Since(msg.Trace.InsertedOnDate).Milliseconds()
	msg.Time = time.Now()

	msg.FetchMetadataFromContext(r.Context)
	msg.Send("http")
}

func Logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ReadAll deletes the content of the body
		bytes, _ := ioutil.ReadAll(r.Body)
		// So we have to restore the io.ReadCloser to its original state
		r.Body = ioutil.NopCloser(bytes2.NewBuffer(bytes))

		var splunkCtx splunk.SplunkContext
		log.SetOutput(&splunkCtx.Buf)
		ctx := context.WithValue(r.Context(), "splunkContext", &splunkCtx)

		request := string(bytes)
		wrapper := ResponseWrapper{
			ResponseWriter: w,
			RequestLog: &RequestLog{
				Request: &request,
			},
			Context: ctx,
		}

		splunk := splunk.HttpLog{}
		splunk.InitRequestLog(r)

		defer wrapper.ExecuteAfterResponding(&splunk)
		next.ServeHTTP(wrapper, r.WithContext(ctx))
	})
}
