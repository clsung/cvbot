FROM golang:1.26-alpine AS builder

RUN apk update && apk upgrade && \
	apk add --no-cache bash \
	bind-tools \
	git

# Install openfortivpn
RUN apk --update upgrade \
    && apk add ca-certificates wget \
    && update-ca-certificates

WORKDIR /go/src/github.com/clsung/cvbot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# simple for simple bot, telsimple add trace for stackdriver logging
# WORKDIR /go/src/github.com/clsung/cvbot/cmd/simple/
# WORKDIR /go/src/github.com/clsung/cvbot/cmd/telsimple/
# tracesim adds trace to stackdriver
#WORKDIR /go/src/github.com/clsung/cvbot/cmd/tracesim/
# logrusimple add stackdriver compatible log
WORKDIR /go/src/github.com/clsung/cvbot/cmd/logrusimple
RUN go build -o /go/bin/cvbot


FROM alpine:3.9

RUN apk --no-cache add ca-certificates openssl ppp curl su-exec bash && rm -rf /var/cache/apk/*;

WORKDIR /
COPY --from=builder /go/bin/cvbot /usr/bin/cvbot

CMD ["/usr/bin/cvbot"]
