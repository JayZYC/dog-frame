package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"sync"
	"time"
)

var (
	// 单例模式
	instance *PrometheusMetrics
	once     sync.Once
)

// PrometheusMetrics 封装 Prometheus 指标
type PrometheusMetrics struct {
	// HTTP 请求计数器
	RequestCounter *prometheus.CounterVec
	// HTTP 请求延迟直方图
	RequestDuration *prometheus.HistogramVec
	// HTTP 请求大小直方图
	RequestSize *prometheus.HistogramVec
	// HTTP 响应大小直方图
	ResponseSize *prometheus.HistogramVec
	// 请求错误计数器
	ErrorCounter *prometheus.CounterVec
}

// GetMetrics 获取 Prometheus 指标实例
func GetMetrics() *PrometheusMetrics {
	once.Do(func() {
		instance = newPrometheusMetrics()
	})
	return instance
}

// newPrometheusMetrics 创建新的 Prometheus 指标
func newPrometheusMetrics() *PrometheusMetrics {
	// 定义标签
	labels := []string{"method", "path", "status"}

	return &PrometheusMetrics{
		// HTTP 请求计数器
		RequestCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "HTTP 请求总数",
			},
			labels,
		),

		// HTTP 请求延迟直方图
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP 请求处理时间（秒）",
				Buckets: prometheus.DefBuckets,
			},
			labels,
		),

		// HTTP 请求大小直方图
		RequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP 请求大小（字节）",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			labels,
		),

		// HTTP 响应大小直方图
		ResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP 响应大小（字节）",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			labels,
		),

		// 请求错误计数器
		ErrorCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_errors_total",
				Help: "HTTP 请求错误总数",
			},
			[]string{"method", "path", "error_type"},
		),
	}
}

// RecordMetrics 记录请求指标
func (p *PrometheusMetrics) RecordMetrics(method, path string, status int, duration time.Duration, requestSize, responseSize int) {
	// 转换状态码为字符串
	statusStr := strconv.Itoa(status)

	// 记录请求计数
	p.RequestCounter.WithLabelValues(method, path, statusStr).Inc()

	// 记录请求处理时间
	p.RequestDuration.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())

	// 记录请求大小
	p.RequestSize.WithLabelValues(method, path, statusStr).Observe(float64(requestSize))

	// 记录响应大小
	p.ResponseSize.WithLabelValues(method, path, statusStr).Observe(float64(responseSize))
}

// RecordError 记录错误
func (p *PrometheusMetrics) RecordError(method, path, errorType string) {
	p.ErrorCounter.WithLabelValues(method, path, errorType).Inc()
}
