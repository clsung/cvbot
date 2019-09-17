package cvbot

import (
	"context"
	"net/http"
	"time"

	teltech "github.com/clsung/logger"
	"github.com/sirupsen/logrus"
)

// Service describes a service that adds things together.
type Service interface {
	Webhook(context.Context, interface{}) (int, error)
	ParseRequest(context.Context, *http.Request) (interface{}, error)
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func LoggingMiddleware(logger Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type Logger interface {
	Printf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type loggingMiddleware struct {
	logger Logger
	next   Service
}

func (mw loggingMiddleware) ParseRequest(ctx context.Context, r *http.Request) (v interface{}, err error) {
	defer func(begin time.Time) {
		mw.logger.Printf("events: %v, time: %v, error: %v", v, time.Since(begin), err)
	}(time.Now())
	return mw.next.ParseRequest(ctx, r)
}

func (mw loggingMiddleware) Webhook(ctx context.Context, r interface{}) (v int, err error) {
	defer func(begin time.Time) {
		if log, ok := mw.logger.(*logrus.Logger); ok {
			e := log.WithFields(logrus.Fields{
				"result": v,
				"error":  err,
				"took":   time.Since(begin),
			})
			if err == nil {
				e.Info("Processed")
			} else {
				e.Error("Failed")
			}
		} else if log, ok := mw.logger.(*teltech.Log); ok {
			if err == nil {
				log.With(teltech.Fields{"result": v, "error": err, "elapsed": time.Since(begin)}).Info("Processed")
			} else {
				log.With(teltech.Fields{"result": v, "error": err, "elapsed": time.Since(begin)}).Error("Failed")
			}
		} else {
			mw.logger.Printf("result: %v, time: %v, error: %v", v, time.Since(begin), err)
		}
	}(time.Now())
	return mw.next.Webhook(ctx, r)
}
