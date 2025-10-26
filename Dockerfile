FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/* ./cmd/
COPY internal/* ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./lambda ./cmd/lambda 
RUN chmod +x bin/lambda


FROM golang:1.25.3-alpine3.22

WORKDIR /app

COPY --from=builder /app/lambda ./lambda

CMD [ "./lambda" ]
