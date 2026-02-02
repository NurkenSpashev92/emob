# ---- Build stage ----
FROM golang:1.25-alpine AS build

# Устанавливаем необходимые утилиты
RUN apk update && apk add --no-cache git bash build-base postgresql-client curl

# Устанавливаем migrate в build-stage
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src .

# Собираем Go приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# Копируем entrypoint
COPY devops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

CMD ["./app"]
