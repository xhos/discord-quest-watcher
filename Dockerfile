FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o orbs

FROM alpine:3.22.1
RUN apk add --no-cache chromium
COPY --from=builder /app/orbs .
CMD ["./orbs"]
