package abstraction

type MetricsObserver interface {
	RecordDuration(metricName string, duration float64, tags map[string]interface{})
	Count(metricName string, value int64, tags map[string]interface{})
	Gauge(metricName string, value float64, tags map[string]interface{})
	Incr(metricName string, tags map[string]interface{})
	Close() error
}
