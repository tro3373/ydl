FROM golang:1.18.0-alpine3.15

WORKDIR /go/src

# COPY . .

RUN set -euxo pipefail \
  && apk update \
  && apk add --no-cache git curl \
  && go install github.com/cosmtrek/air@latest

CMD air
