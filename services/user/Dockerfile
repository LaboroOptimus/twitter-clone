FROM golang:1.24-alpine

WORKDIR /app
RUN apk add --no-cache git

# копируем весь проект сразу
COPY . .

WORKDIR /app/services/user
RUN go mod download
RUN go build -o user .

EXPOSE 8080
CMD ["./user"]
