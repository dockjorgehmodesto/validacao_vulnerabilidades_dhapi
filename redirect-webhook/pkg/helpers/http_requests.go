package helpers

import (
	"bytes"
	"context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

type HttpResponse struct {
	Body                 string `json:"body"`
	Status               int    `json:"status"`
	Duration             int64  `json:"duration"`
	TotalBytesThroughput int64  `json:"totalBytesThroughput"`
}

func HttpRequest(ctx context.Context, client *http.Client, url string, headers map[string]string, body []byte, method string) (*HttpResponse, error) {
	tracerSpan, ctx := tracer.StartSpanFromContext(ctx, "httprequest")
	defer tracerSpan.Finish()
	var requestTotalBytes int64
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if out, err := httputil.DumpRequestOut(request, true); err == nil {
		requestTotalBytes += int64(len(out))
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["SendDate"] = time.Now().Format("2006-01-02T15:04:05.000Z")

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	timer := time.Now()
	response, err := client.Do(request)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	duration := time.Since(timer).Milliseconds()

	if out, err := httputil.DumpResponse(response, true); err == nil {
		requestTotalBytes += int64(len(out))
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Body:                 string(responseBody),
		Status:               response.StatusCode,
		Duration:             duration,
		TotalBytesThroughput: requestTotalBytes,
	}, nil

}
