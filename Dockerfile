FROM golang:1.25.0-alpine AS builder

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bookmark-management ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /opt/app/bookmark-management /app/bookmark-management
COPY --from=builder /opt/app/docs /app/docs


CMD ["./bookmark-management"]