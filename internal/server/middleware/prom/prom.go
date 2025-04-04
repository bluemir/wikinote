package prom

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bluemir/wikinote/internal/buildinfo"
)

var (
	metricsRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: buildinfo.AppName + "_request_total",
			Help: "web server request count",
		},
		[]string{"method", "url", "code"},
	)
	metricsRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: buildinfo.AppName + "_request_duration",
			Help: "the time server took to handle the request.",
			Buckets: []float64{
				float64(1 * time.Millisecond),
				float64(10 * time.Millisecond),
				float64(100 * time.Millisecond),
				float64(500 * time.Millisecond),
				float64(1 * time.Second),
				float64(3 * time.Second),
				float64(5 * time.Second),
			},
		},
		[]string{"method", "url", "code"},
	)
)

func init() {
	prometheus.MustRegister(metricsRequestCount)
	prometheus.MustRegister(metricsRequestDuration)
}
func Metrics(labels ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		/* if need sampling
		if rand.Intn(100) < 1 {
			return // skip & next
		}
		*/

		start := time.Now()
		c.Next()
		end := time.Now()

		label := prometheus.Labels(map[string]string{
			"method": c.Request.Method,
			"url":    c.FullPath(),
			"code":   strconv.Itoa(c.Writer.Status()),
		})

		metricsRequestCount.With(label).Inc()
		metricsRequestDuration.With(label).Observe(float64(end.Sub(start).Milliseconds()))
	}
}
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
