package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogRequest is a gin middleware that logs useful informations on each request
// as they come.
//
// If the request's context contains a logger, it will be used. Otherwise, the
// specified fallback will be used instead.
func LogRequest(fallbackLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := fallbackLogger

		// try extract logger from context
		sugaredLogger := GetLoggerFromContext(c)
		if sugaredLogger != nil {
			logger = sugaredLogger.Desugar()
		}

		// execute the request
		start := time.Now()
		c.Next()

		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// log request informations
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", start.Format(time.RFC3339)),
		)
	}
}

func GetLoggerFromContext(c *gin.Context) *zap.SugaredLogger {
	value, ok := c.Get("logger")
	if !ok {
		return nil
	}

	l, ok := value.(*zap.SugaredLogger)
	if !ok {
		panic("logger in context is not a *zap.SugaredLogger")
	}

	return l
}
