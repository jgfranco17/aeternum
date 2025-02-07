package system

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpLastRequestReceivedTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_last_request_received_time",
			Help: "Time when the last request was processed",
		}, []string{"path", "method"},
	)
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		HttpLastRequestReceivedTime.WithLabelValues(c.FullPath(), c.Request.Method).SetToCurrentTime()
	}
}
