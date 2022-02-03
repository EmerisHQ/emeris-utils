package logging

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ctxKey string

const (
	CorrelationIDName         ctxKey = "correlation_id"
	IntCorrelationIDName      ctxKey = "int_correlation_id"
	ExternalCorrelationIDName string = "X-Correlation-Id"
)

// CorrelationIDMiddleware adds correlationID if it's not specified in HTTP request
func CorrelationIDMiddleware(l *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		addCorrelationID(c, l)
	}
}

func addCorrelationID(c *gin.Context, l *zap.SugaredLogger) {
	ctx := c.Request.Context()

	correlationID := c.Request.Header.Get(ExternalCorrelationIDName)

	if correlationID != "" {
		ctx = context.WithValue(ctx, CorrelationIDName, correlationID)
		c.Writer.Header().Set(ExternalCorrelationIDName, correlationID)
		l = l.With(CorrelationIDName, correlationID)
	}

	id, _ := uuid.NewV4()

	ctx = context.WithValue(ctx, IntCorrelationIDName, id.String())
	l = l.With(IntCorrelationIDName, id)

	c.Set("logger", l)

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
