# Этап сборки
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app_aircraft server.go

# Финальный этап
FROM alpine:latest
WORKDIR /root/

# Копируем сам бинарник из этапа builder
COPY --from=builder /app/app_aircraft .

# Копируем файл настроек из этапа builder
COPY --from=builder /app/config.yml .

EXPOSE 9081
ENTRYPOINT ["./app_aircraft"]
