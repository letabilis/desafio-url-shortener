FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o solucao-url-shortener ./cmd


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/solucao-url-shortener .
EXPOSE 8080
CMD ["./solucao-url-shortener"]



