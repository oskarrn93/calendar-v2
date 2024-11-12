FROM golang:1.23.3-alpine3.20

WORKDIR /asset

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bootstrap

CMD [ "./bootstrap" ]
