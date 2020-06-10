package cvbot

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/clsung/cvbot/cv/aws"
	"github.com/clsung/cvbot/cv/gcp"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// CVApp app
type CVApp struct {
	bot *linebot.Client
}

// NewCVApp returns a naïve, stateless implementation of Service.
func NewCVApp() Service {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &CVApp{
		bot: bot,
	}

}

// ParseRequest parse http.Request into events
func (app *CVApp) ParseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return app.bot.ParseRequest(r)
}

// Webhook handles parsed events
func (app *CVApp) Webhook(ctx context.Context, r interface{}) (int, error) {
	events := r.([]*linebot.Event)
	for _, event := range events {
		replyMsgs, err := app.HandleEvent(ctx, event)
		if err != nil {
			return 0, errors.Wrap(err, "handle event failed")
		}
		if _, err := app.bot.ReplyMessage(
			event.ReplyToken,
			replyMsgs...,
		).Do(); err != nil {
			return 0, errors.Wrap(err, "reply message failed")
		}
	}
	return len(events), nil
}

// HandleEvent implements event handler with provider-based logic
func (app *CVApp) HandleEvent(ctx context.Context, event *linebot.Event) (replyMsgs []linebot.SendingMessage, err error) {
	ctx, span := trace.StartSpan(ctx, "handle_event")
	defer span.End()
	if event.Type == linebot.EventTypeMessage {
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			span.AddAttributes(trace.StringAttribute("type", "text"))
			//replyMsgs = append(replyMsgs, linebot.NewTextMessage(message.Text))
		case *linebot.ImageMessage:
			span.AddAttributes(trace.StringAttribute("type", "image"))
			var retMsgs []linebot.SendingMessage
			if retMsgs, err = app.handleImageAsync(ctx, message); err != nil {
				return replyMsgs, errors.Wrap(err, "handle image failed")
			}
			replyMsgs = append(replyMsgs, retMsgs...)
		case *linebot.StickerMessage:
			span.AddAttributes(trace.StringAttribute("type", "sticker"))
			return replyMsgs, errors.New("STICKER ERROR")
		default:
			span.AddAttributes(trace.StringAttribute("type", "unknown"))
			replyMsgs = append(replyMsgs, linebot.NewTextMessage("not supported yet"))
		}
	}
	return replyMsgs, nil
}

func (app *CVApp) handleImage(ctx context.Context, message *linebot.ImageMessage) (replyMsgs []linebot.SendingMessage, err error) {
	ctx, span := trace.StartSpan(ctx, "handle_image")
	defer span.End()
	_, err = app.handleFaceRecognition(ctx, message.ID, func(gcpBuf io.Reader, awsBuf io.Reader) error {
		_, span := trace.StartSpan(ctx, "aws")
		awsFaces, err := aws.FaceDetect(awsBuf)
		span.End()
		if err != nil {
			return errors.Wrap(err, "aws face detect error")
		}
		replyMsgs = append(replyMsgs, linebot.NewTextMessage(string(awsFaces)))
		_, span = trace.StartSpan(ctx, "gcp")
		gcpFaces, err := gcp.FaceDetect(gcpBuf)
		span.End()
		if err != nil {
			return errors.Wrap(err, "gcp face detect error")
		}
		replyMsgs = append(replyMsgs, linebot.NewTextMessage(string(gcpFaces)))
		return nil
	})
	return replyMsgs, err
}

func (app *CVApp) handleFaceRecognition(ctx context.Context, messageID string, callback func(io.Reader, io.Reader) error) ([]linebot.SendingMessage, error) {
	_, span := trace.StartSpan(ctx, "handle_facereg")
	defer span.End()
	_, span2 := trace.StartSpan(ctx, "get_image")
	content, err := app.bot.GetMessageContent(messageID).Do()
	span2.End()
	if err != nil {
		return nil, errors.Wrap(err, "get image from line error")
	}
	defer content.Content.Close()
	var gcpBuf, awsBuf bytes.Buffer
	w := io.MultiWriter(&gcpBuf, &awsBuf)

	if _, err := io.Copy(w, content.Content); err != nil {
		return nil, errors.Wrap(err, "multiwriter copy error")
	}
	return nil, callback(io.Reader(&gcpBuf), io.Reader(&awsBuf))
}

func (app *CVApp) handleImageAsync(ctx context.Context, message *linebot.ImageMessage) (replyMsgs []linebot.SendingMessage, err error) {
	ctx, span := trace.StartSpan(ctx, "handle_image_async")
	defer span.End()
	_, err = app.handleFaceRecognition(ctx, message.ID, func(gcpBuf io.Reader, awsBuf io.Reader) error {
		dataCh := make(chan string, 1)
		go func() {
			_, span := trace.StartSpan(ctx, "aws")
			awsFaces, _ := aws.FaceDetect(awsBuf)
			span.End()
			dataCh <- string(awsFaces)
			return
		}()
		go func() {
			_, span := trace.StartSpan(ctx, "gcp")
			// uncomment it for demo purpose
			// time.Sleep(time.Duration(rand.Int31n(2000)) * time.Millisecond)
			gcpFaces, _ := gcp.FaceDetect(gcpBuf)
			span.End()
			dataCh <- string(gcpFaces)
			return
		}()
		ret := <-dataCh
		replyMsgs = append(replyMsgs, linebot.NewTextMessage(string(ret)))
		return nil
	})
	return replyMsgs, err
}
