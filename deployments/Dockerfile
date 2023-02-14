FROM node:18.12.0-alpine as frontend-builder

WORKDIR /app

COPY internal/server/web .

RUN npm i -g pnpm && pnpm install && pnpm run build

FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

COPY --from=frontend-builder /app/dist /app/internal/server/web

RUN go build -ldflags="-s -w" -o beaver_server ./cmd/beaver_server

FROM alpine:3.17.1  

WORKDIR /app

RUN apk add libc6-compat=1.2.3-r4 --no-cache && rm -rf /var/cache/apk/*

COPY --from=builder /app/beaver_server /app/

VOLUME [ "/app/config" ]
VOLUME [ "/app/data" ]

ENTRYPOINT ["./beaver_server"]
CMD [ "--config", "/app/config/beaver_server.yaml" ]
