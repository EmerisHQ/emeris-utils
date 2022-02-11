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
		ctx = context.WithValue(ctx, string(CorrelationIDName), correlationID)
		c.Writer.Header().Set(ExternalCorrelationIDName, correlationID)
		l = l.With(string(CorrelationIDName), correlationID)
	}

	id, err := uuid.NewV4()
	if err != nil {
		l.Errorf("Error while creating new internal correlation id error: %w", err)
	}

	ctx = context.WithValue(ctx, string(IntCorrelationIDName), id.String())
	l = l.With(string(IntCorrelationIDName), id)

	c.Set("logger", l)

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
