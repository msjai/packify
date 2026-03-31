FROM golang:1.26.1-alpine AS builder

WORKDIR /app
RUN apk add --update gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o ./bin/packify ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/packify .
COPY --from=builder /app/config ./config

CMD ["./packify", "--config=./config/config.yaml"]