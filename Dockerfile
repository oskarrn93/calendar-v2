ARG GO_VERSION=1.27.0

FROM golang:${GO_VERSION}-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/lambda ./cmd/lambda

FROM alpine:3.24

RUN apk add --no-cache ca-certificates

WORKDIR /asset

COPY --from=builder /app/bin/lambda ./bootstrap

CMD [ "./bootstrap" ]
