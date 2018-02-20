FROM golang:1.9.2-alpine3.7
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
RUN apk add --no-cache git build-base \
    && go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/authserver
WORKDIR /go/src/github.com/authserver
COPY Gopkg.* main.go ./
ADD commands commands
ADD middlewares middlewares
ADD oauth2 oauth2
ADD user user
RUN dep ensure \
    && go build -o app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/authserver/app /usr/local/bin/
