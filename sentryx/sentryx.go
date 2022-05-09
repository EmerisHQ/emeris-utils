package sentryx

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// GinMiddleware returns a middlewares that properly logs the query and its
// error in Sentry.
func GinMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		ctx = sentry.SetHubOnContext(ctx, hub)
	}
	c.Request = c.Request.WithContext(ctx)
	hub.Scope().SetRequest(c.Request)

	// start span for the request
	span, ctx := StartSpan(c.Request.Context(), "http.server",
		sentry.TransactionName(fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())),
		sentry.ContinueFromRequest(c.Request),
	)
	for _, param := range c.Params {
		span.SetTag("param."+param.Key, param.Value)
	}
	defer func() {
		span.Finish()
		if err := recover(); err != nil {
			// Send error to sentry
			hub.RecoverWithContext(ctx, err)
			// repanic
			panic(err)
		}
	}()
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// StartSpan calls sentry.StartSpan and returns the span and the span context.
// This is handy to ensure correct hierarchy of span if other spans are started
// after this one. Indeed you need to use the span.Context() to start a child
// span. Recommanded usage is:
//
//     span, ctx := sentryx.StartSpan(ctx, "operation")
//     defer span.Finish()
func StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) (*sentry.Span, context.Context) {
	if ginctx, ok := ctx.(*gin.Context); ok {
		// Ensure we use the good context (released version of gin doesn't fallback
		// to the request context).
		ctx = ginctx.Request.Context()
	}
	s := sentry.StartSpan(ctx, operation, options...)
	return s, s.Context()
}
