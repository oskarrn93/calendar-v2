FROM golang:1.26rc2-alpine3.22 AS builder

RUN apk update && apk add --no-cache make git

WORKDIR /app

RUN mkdir ./bin

COPY Makefile ./
COPY go.mod ./
COPY go.sum ./

RUN make install

COPY . ./

RUN make build

FROM golang:1.26rc2-alpine3.22

WORKDIR /asset

COPY --from=builder /app/bin/lambda ./bootstrap

CMD [ "./bootstrap" ]
