package logging

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LogRequest is a gin middleware that logs useful informations on each request as they come.
func LogRequest(c *gin.Context) {
	l, err := GetLoggerFromContext(c)
	if err != nil {
		panic(fmt.Errorf("Can't get logger from context error:%w", err))
	}
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

func GetLoggerFromContext(c *gin.Context) (*zap.SugaredLogger, error) {
	value, ok := c.Get("logger")
	if !ok {
		return nil, fmt.Errorf("logger does not exists in context")
	}

	l, ok := value.(*zap.SugaredLogger)
	if !ok {
		return nil, fmt.Errorf("invalid logger format in context")
	}

	return l, nil
}
