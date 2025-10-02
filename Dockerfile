FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o discord-quest-watcher ./cmd

FROM alpine:3.22.1
RUN apk add --no-cache chromium
COPY --from=builder /app/discord-quest-watcher .
CMD ["./discord-quest-watcher"]
