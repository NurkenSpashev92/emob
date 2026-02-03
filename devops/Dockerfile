FROM golang:1.25-alpine AS builder

RUN apk update && apk add --no-cache git bash build-base

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY .env /app/.env

COPY src .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

COPY devops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]

CMD ["/app/main"]
