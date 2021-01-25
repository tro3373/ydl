FROM golang:1.15.6-alpine3.12

WORKDIR /go/src

COPY . .

RUN set -euxo pipefail \
  && apk update \
  && apk add --no-cache git curl \
  && go get -u github.com/cosmtrek/air

CMD air
