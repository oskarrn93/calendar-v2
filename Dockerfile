FROM golang:1.27.1-alpine3.24 AS builder

RUN apk update && apk add --no-cache make git

WORKDIR /app

RUN mkdir ./bin

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN make install

COPY . ./

RUN make build

FROM golang:1.27.1-alpine3.24

WORKDIR /asset

COPY --from=builder /app/bin/lambda ./bootstrap

CMD [ "./bootstrap" ]
