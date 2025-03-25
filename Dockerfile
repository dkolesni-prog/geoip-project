FROM golang:1.23.7-alpine AS builder
LABEL authors="daniilkolesnik"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/whois/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/main .

COPY --from=frontend-builder /app/build ./frontend/build

EXPOSE 8080
CMD ["./main"]
