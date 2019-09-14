package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/clsung/cvbot"
	"github.com/clsung/logger"
	"github.com/line/line-bot-sdk-go/linebot"

	"go.opencensus.io/trace"
)

func main() {
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	sd, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "alansandbox",
		// MetricPrefix helps uniquely identify your metrics.
		MetricPrefix: "tracesim",
	})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	// Register it as a trace exporter
	trace.RegisterExporter(sd)

	app := cvbot.LoggingMiddleware(logger.New())(cvbot.NewCVApp())
	app = cvbot.TraceMiddleware()(app)
	// setup trace config, production please use trace.ProbabilitySampler
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "webhook")
		defer span.End()
		events, err := app.ParseRequest(ctx, req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		app.Webhook(ctx, events)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte(`{"health": "OK"}`))
		if err != nil {
			log.Fatal(err)
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
