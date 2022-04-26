package sentryx

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// StartSpan calls sentry.StartSpan and returns the span and the span context.
// This is handy to ensure correct hierarchy of span if other spans are started
// after this one. Indeed you need to use the span.Context() to start a child
// span. Recommanded usage is:
//
//     span, ctx := sentryx.StartSpan(ctx, "operation")
//     defer span.Finish()
func StartSpan(ctx context.Context, operation string, options ...sentry.SpanOption) (*sentry.Span, context.Context) {
	s := sentry.StartSpan(ctx, operation, options...)
	return s, s.Context()
}
