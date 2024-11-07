package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()

		logger.WithFields(logger.Fields{
			"status_code":  statusCode,
			"latency_time": latency,
			"client_ip":    c.ClientIP(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
		}).Info("Request details")
	}
}
