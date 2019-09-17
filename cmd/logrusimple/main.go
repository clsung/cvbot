package main

import (
	"fmt"
	"io/ioutil"
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
		},
	})

	// create hook using service account credentials
	if metadata.OnGCE() {
		// stackdriver query:
		//  resource.type = "cloud_run_revision"
		//  resource.labels.service_name = "cvbot"
		//  resource.labels.location = "us-central1"
		//  NOT logName: "cloudaudit.googleapis.com"
		//  NOT (logName : "varlog%2Fsystem" AND severity = DEBUG)
		//  severity>=DEFAULT
		ProjectID, _ := metadata.ProjectID()
		Zone, _ := metadata.Zone() // zone, such as "us-central1-b".
		h, err := sdhook.New(
			sdhook.GoogleComputeCredentials(""), // use default service account
			sdhook.ProjectID(ProjectID),
			sdhook.LogName(os.Getenv("K_SERVICE")), // must call after ProjectID()
			sdhook.Resource(sdhook.ResTypeCloudRunRevision, map[string]string{
				"project_id":         ProjectID,
				"service_name":       os.Getenv("K_SERVICE"),
				"revision_name":      os.Getenv("K_REVISION"),
				"location":           Zone[:len(Zone)-2], // location, such as "us-central1"
				"configuration_name": os.Getenv("K_CONFIGURATION"),
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		logger.Hooks.Add(h)
		logger.Out = ioutil.Discard // don't print to stderr, omit it
	}
	app := cvbot.LoggingMiddleware(logger)(cvbot.NewCVApp())
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		events, err := app.ParseRequest(req.Context(), req)
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
