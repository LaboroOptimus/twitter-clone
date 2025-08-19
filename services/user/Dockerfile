# Используем официальный образ Go
FROM golang:1.24-alpine

# Установим рабочую директорию внутри контейнера
WORKDIR /app

# Установим зависимости для компиляции Go и GORM с Postgres
RUN apk add --no-cache git

COPY go.* ./
RUN if [ -f go.sum ]; then go mod download; fi

COPY . .
WORKDIR /app/services/user
RUN go build -o user .
EXPOSE 8080
CMD ["./user"]
