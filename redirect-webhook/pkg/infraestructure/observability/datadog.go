package observability

import (
	"fmt"
	"log"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"github.com/DataDog/datadog-go/statsd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type DatadogObserver struct {
	client    statsd.ClientInterface
	ddEnabled bool
}

func NewDatadogObserver() *DatadogObserver {
	ddEnabled := viper.GetBool("datadog.enabled")
	if ddEnabled {
		startDataDogProfilerAndTracer()
		addr := fmt.Sprintf("%s:%s", global.Viper.GetString("DATADOG_AGENT_HOST"), global.Viper.GetString("DATADOG_AGENT_PORT"))
		statsdClient, err := statsd.New(addr, statsd.WithNamespace(global.Viper.GetString("app.serviceName")))
		if err != nil {
			log.Fatalf("erro ao iniciar statsd: %v", err)
		}
		return &DatadogObserver{client: statsdClient, ddEnabled: ddEnabled}
	}
	return &DatadogObserver{ddEnabled: ddEnabled}
}

// RecordDuration usa Distribution, que é o ideal para gerar percentis globais (p50, p90, p95, p99)
// com alta precisão na interface do Datadog para durações em milissegundos.
func (d *DatadogObserver) RecordDuration(metricName string, duration float64, tags map[string]interface{}) {
	if d.ddEnabled {
		_ = d.client.Distribution(metricName, duration, buildTags(tags), 1.0)
	}
}
func (d *DatadogObserver) Incr(metricName string, tags map[string]interface{}) {
	if d.ddEnabled {
		_ = d.client.Incr(metricName, buildTags(tags), 1.0)
	}
}

// Count adiciona um valor arbitrário a um contador. Ideal para somar bytes processados, itens em um lote, etc.
func (d *DatadogObserver) Count(metricName string, value int64, tags map[string]interface{}) {
	if d.ddEnabled {
		_ = d.client.Count(metricName, value, buildTags(tags), 1.0)
	}
}

// Gauge mede o valor de uma métrica em um instante específico (ex: uso de CPU, conexões ativas, tamanho de fila).
func (d *DatadogObserver) Gauge(metricName string, value float64, tags map[string]interface{}) {
	if d.ddEnabled {
		_ = d.client.Gauge(metricName, value, buildTags(tags), 1.0)
	}
}
func (d DatadogObserver) Close() error {
	if d.ddEnabled {
		tracer.Stop()
		profiler.Stop()
		return d.client.Close()
	}
	return nil
}
func buildTags(tags map[string]interface{}) []string {
	if len(tags) == 0 {
		return nil
	}
	ddTags := make([]string, 0, len(tags))
	for k, v := range tags {
		ddTags = append(ddTags, fmt.Sprintf("%s:%v", k, v))
	}
	return ddTags
}
func startDataDogProfilerAndTracer() {
	err := profiler.Start(
		profiler.WithService(global.Viper.GetString("app.serviceName")),
		profiler.WithEnv(global.Viper.GetString("app.profile")),
		profiler.WithVersion(global.VersionInfo.Current),
		profiler.WithTags(fmt.Sprintf("pod_name:%s", global.Viper.GetString("pod.name"))),
	)
	if err != nil {
		logging.Logger.Error("datadog profiler error", zap.Error(err))
	}
	tracer.Start(
		tracer.WithService(global.Viper.GetString("app.serviceName")),
		tracer.WithEnv(global.Viper.GetString("app.profile")),
		tracer.WithRuntimeMetrics(),
		tracer.WithServiceVersion(global.VersionInfo.Current),
		tracer.WithGlobalTag("pod_name", global.Viper.GetString("pod.name")),
	)
}
