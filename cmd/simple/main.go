package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/clsung/cvbot"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	app := cvbot.NewCVApp()
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
		v, err := app.Webhook(req.Context(), events)
		log.Printf("events: %d, err: %v", v, err)
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
