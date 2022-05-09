package ginsentry

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// Recover can be called while recovering from a panic, pass the gin.Context and
// the returned value of recover() to log this event to Sentry.
func Recover(c *gin.Context, recoverVal interface{}) {
	hub := sentry.GetHubFromContext(c.Request.Context())
	hub.RecoverWithContext(
		context.WithValue(c.Request.Context(), sentry.RequestContextKey, c.Request),
		recoverVal,
	)
}

// Middleware is a gin middleware that logs requests to Sentry.
func Middleware(ctx *gin.Context) {
	// set Hub inside Request context
	reqCtx := ctx.Request.Context()
	hub := sentry.GetHubFromContext(reqCtx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
		reqCtx = sentry.SetHubOnContext(reqCtx, hub)
	}
	ctx.Request = ctx.Request.WithContext(reqCtx)
	hub.Scope().SetRequest(ctx.Request)

	// start span for the request
	span := sentry.StartSpan(ctx.Request.Context(), "http.server",
		sentry.TransactionName(fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())),
		sentry.ContinueFromRequest(ctx.Request),
	)
	defer span.Finish()
	for _, param := range ctx.Params {
		span.SetTag(param.Key, param.Value)
	}
	ctx.Request = ctx.Request.WithContext(span.Context())

	ctx.Next()
}

// StartSpan is a shorthand for sentry.StartSpan(c.Request.Context(), ...)
func StartSpan(c *gin.Context, operation string, options ...sentry.SpanOption) *sentry.Span {
	s := sentry.StartSpan(c.Request.Context(), operation, options...)
	return s
}
