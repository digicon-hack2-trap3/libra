FROM golang:1.21-alpine
WORKDIR /app

COPY go.mod go.sum *.go ./

RUN go mod download

RUN go build -o app .
ENTRYPOINT ./app