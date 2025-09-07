FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

# Копируем всё репо сразу — пути ../../pkg/* уже существуют
COPY . .

# Теперь go mod download видит ../../pkg/redisx и ../../pkg/authjwt
WORKDIR /app/services/post
RUN go mod download
RUN go build -o /bin/posts-service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/posts-service ./posts-service
EXPOSE 8081
CMD ["./posts-service"]
