FROM golang:1.24-alpine
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# модули (go.sum может отсутствовать)
COPY go.* ./
RUN if [ -f go.sum ]; then go mod download; fi

# исходники
COPY . .

# сборка из корня
RUN go build -o posts .

EXPOSE 8081
CMD ["./posts"]
