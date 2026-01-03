package middleware

import (
	"codim/pkg/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogger returns a gin.HandlerFunc (middleware) that logs requests using the custom logger
func GinLogger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)

		entry := log.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"method":  c.Request.Method,
			"path":    path,
			"latency": latency,
			"error":   c.Errors.ByType(gin.ErrorTypePrivate).String(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			if c.Writer.Status() >= 500 {
				entry.Error("Request failed")
			} else if c.Writer.Status() >= 400 {
				entry.Warn("Request failed")
			} else {
				entry.Info("Request processed")
			}
		}
	}
}
