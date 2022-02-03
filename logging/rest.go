package logging

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var l *zap.Logger

// LogRequest is a gin middleware that logs useful informations on each request as they come.
func LogRequest(logger *zap.Logger) gin.HandlerFunc {
	l = logger
	return log
}

func log(c *gin.Context) {
	start := time.Now()

	c.Next()

	// some evil middlewares modify this values
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	l.Info(path,
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("time", start.Format(time.RFC3339)),
	)
}

func GetLoggerFromContext(c *gin.Context) *zap.SugaredLogger {
	value, ok := c.Get("logger")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "logger does not exists in context")
		return nil
	}

	l, ok := value.(*zap.SugaredLogger)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid logger format in context")
		return nil
	}

	return l
}
