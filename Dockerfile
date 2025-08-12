FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o url-shortener ./cmd


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
CMD ["./url-shortener"]



