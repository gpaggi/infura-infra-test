FROM golang:1.16.5-alpine3.13 AS builder
ARG VERSION
WORKDIR /go/app

RUN apk add --no-cache bash git

COPY . /go/app/
RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o bin/ethapi -ldflags "-X github.com/gpaggi/ethapi/version.Version=$(VERSION) -s -w" .

FROM alpine:3.13
RUN adduser -s /sbin/nologin -D -H -u 1000 -g ethapi ethapi

USER ethapi

COPY --from=builder /go/app/bin/ethapi /usr/sbin/ethapi

CMD ["/usr/sbin/ethapi"]