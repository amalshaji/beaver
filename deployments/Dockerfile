FROM node:18.12.1-alpine as frontend-builder

WORKDIR /app

COPY internal/server/web .

RUN npm i -g pnpm && pnpm install && pnpm run build

FROM golang:1.20.1 AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN mkdir -p /app/internal/server/web/dist

COPY --from=frontend-builder /app/dist /app/internal/server/web/dist

RUN CGO_ENABLED=1 go build -ldflags="-s -w -linkmode external -extldflags \"-static\"" -o beaver_server ./cmd/beaver_server

FROM alpine:3.17.2  

WORKDIR /app

RUN apk add libc6-compat=1.2.3-r4 --no-cache && rm -rf /var/cache/apk/*

COPY --from=builder /app/beaver_server /app/

VOLUME [ "/app/config" ]
VOLUME [ "/app/data" ]

ENTRYPOINT ["./beaver_server"]
CMD [ "--config", "/app/config/beaver_server.yaml" ]
