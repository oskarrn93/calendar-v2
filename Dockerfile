FROM golang:1.22.2-alpine3.19

WORKDIR /asset

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bootstrap

CMD [ "./bootstrap" ]
