FROM golang:1.22-alpine AS builder
LABEL authors="daniilkolesnik"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/main .
COPY --from=builder /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./main"]