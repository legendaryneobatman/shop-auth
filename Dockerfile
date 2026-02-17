FROM golang:1.25-alpine

WORKDIR /app

# Устанавливаем Air
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

RUN chown -R 1000:1000 /app

# Запускаем Air, который будет следить за кодом
CMD ["air", "-c", ".air.toml"]