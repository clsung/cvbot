package cvbot

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
)

// TraceMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func TraceMiddleware() Middleware {
	return func(next Service) Service {
		return traceMiddleware{
			next: next,
		}
	}
}

type traceMiddleware struct {
	next Service
}

func (mw traceMiddleware) ParseRequest(ctx context.Context, r *http.Request) (v interface{}, err error) {
	ctx, span := trace.StartSpan(ctx, "parse_request")
	span.AddAttributes(trace.StringAttribute("hello", "world"))
	defer span.End()
	return mw.next.ParseRequest(ctx, r)
}

func (mw traceMiddleware) Webhook(ctx context.Context, r interface{}) (v int, err error) {
	ctx, span := trace.StartSpan(ctx, "callback")
	defer span.End()
	return mw.next.Webhook(ctx, r)
}
