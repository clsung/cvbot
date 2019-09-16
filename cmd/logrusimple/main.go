package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/clsung/cvbot"
	"github.com/knq/sdhook"
	"github.com/sirupsen/logrus"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
	// create hook using service account credentials
	if metadata.OnGCE() {
		log.Printf("test sdhook support")
		h, err := sdhook.New(
			sdhook.GoogleServiceAccountCredentialsFile("/auth.json"),
			//sdhook.GoogleComputeCredentials(""), // use default service account
		)
		log.Printf("test sdhook support %v", err)
		if err != nil {
			log.Fatal(err)
		}
		logger.Hooks.Add(h)
		log.Printf("sdhook support hooked")
	} else {
		log.Printf("no sdhook support")
	}
	app := cvbot.LoggingMiddleware(logger)(cvbot.NewCVApp())
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		events, err := app.ParseRequest(req.Context(), req)
		log.Printf("events: %d, err: %v", events, err)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		app.Webhook(req.Context(), events)
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
