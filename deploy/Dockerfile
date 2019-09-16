FROM golang:1.12-alpine as builder

RUN apk update && apk upgrade && \
	apk add --no-cache bash \
	bind-tools \
	git

# Install openfortivpn
RUN apk --update upgrade \
    && apk add ca-certificates wget \
    && update-ca-certificates

ENV GO111MODULE off
RUN go get -u "github.com/aws/aws-sdk-go/aws" \
    "github.com/aws/aws-sdk-go/aws/session" \
    "github.com/aws/aws-sdk-go/service/rekognition" \
    "cloud.google.com/go/vision/apiv1" \
    "github.com/sirupsen/logrus" \
    "github.com/knq/sdhook" \
	"go.opencensus.io/trace" \
	"contrib.go.opencensus.io/exporter/stackdriver" \
	"github.com/clsung/logger" \
	"github.com/pkg/errors" \
    "github.com/line/line-bot-sdk-go/linebot"
COPY . /go/src/github.com/clsung/cvbot/
# simple for simple bot, telsimple add trace for stackdriver logging
# WORKDIR /go/src/github.com/clsung/cvbot/cmd/simple/
# WORKDIR /go/src/github.com/clsung/cvbot/cmd/telsimple/
# tracesim adds trace to stackdriver
#WORKDIR /go/src/github.com/clsung/cvbot/cmd/tracesim/
# logrusimple add stackdriver compatible log
WORKDIR /go/src/github.com/clsung/cvbot/cmd/logrusimple/
RUN go build -o /go/bin/cvbot


FROM alpine:3.9

RUN apk --no-cache add ca-certificates openssl ppp curl su-exec bash && rm -rf /var/cache/apk/*;

WORKDIR /
COPY --from=builder /go/bin/cvbot /usr/bin/cvbot

CMD ["/usr/bin/cvbot"]
