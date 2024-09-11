FROM golang:latest AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bot ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bot /usr/local/bin/bot
ENTRYPOINT ["/usr/local/bin/bot"]
