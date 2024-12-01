// middleware/logger.go
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger is a Gin middleware that logs HTTP requests using zerolog.
func Logger(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Log the details of the request
		log.Info().
			Int("status", statusCode).
			String("method", c.Request.Method).
			String("path", c.Request.URL.Path).
			Dur("latency", latency).
			Msg("Handled request")
	}
}