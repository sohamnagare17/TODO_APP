package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "todo_http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "route", "status"},
)

var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "todo_http_request_duration_seconds",
		Help: "Duration of HTTP requests",
	},
	[]string{"method", "route"},
)

var HttpErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "todo_http_errors_total",
		Help: "Total number of HTTP errors",
	},
	[]string{"method", "route"},
)

var DBQueryDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "todo_db_query_duration_seconds",
		Help: "Duration of DB queries",
	},
	[]string{"query", "table"},
)

var DBErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "todo_db_errors_total",
		Help: "Total DB errors",
	},
	[]string{"query", "table"},
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpErrorsTotal)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(DBErrorsTotal)
}
