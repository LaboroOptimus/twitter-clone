FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git


COPY . .

# Теперь go mod download видит ../../pkg/redisx и ../../pkg/authjwt
WORKDIR /app/services/notifications
RUN go mod download
RUN go build -o /bin/notifications-service .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/notifications-service ./notifications-service
EXPOSE 8082
CMD ["./notifications-service"]
