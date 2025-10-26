FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

RUN apk update && apk add --no-cache make

COPY go.mod go.sum ./
COPY Makefile ./

RUN make install

COPY cmd/* ./cmd/
COPY internal/* ./internal/

RUN mkdir -p ./bin

RUN make build


FROM golang:1.25.3-alpine3.22

WORKDIR /app

COPY --from=builder /app/bin/lambda ./lambda

CMD [ "./lambda" ]
