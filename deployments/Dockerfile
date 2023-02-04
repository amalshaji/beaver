FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -ldflags="-s -w" -o beaver_server ./cmd/beaver_server

FROM alpine:3.17.1  

WORKDIR /app

RUN apk add libc6-compat=1.2.3-r4 --no-cache && rm -rf /var/cache/apk/*

COPY --from=builder /app/beaver_server /app/

ENTRYPOINT ["./beaver_server"]
