APP_NAME=app
COMPOSE=docker compose
ENV_FILE=.env

.PHONY: help up down build restart logs ps exec app postgres create_migration migrations_up clean prune

help:
	@echo ""
	@echo "Available commands:"
	@echo "  make install          🚀 Deploy project"
	@echo "  make up               🚀 Start containers"
	@echo "  make down             🛑 Stop containers"
	@echo "  make build            � Build containers"
	@echo "  make restart          � Restart containers"
	@echo "  make logs             📜 Show logs"
	@echo "  make ps               📦 Show containers"
	@echo "  make app              🐹 Enter app container"
	@echo "  make postgres         🐘 Enter postgres container"
	@echo "  make create_migration 📝 Create migrations"
	@echo "  make migrations_up    ⬆️ Run migrations"
	@echo "  make swagger          📖 Generate Swagger docs"
	@echo "  make clean            🧹 Remove containers + volumes"
	@echo "  make prune            💣 Docker system prune"
	@echo ""

## -----------------------------
## 🐳 Docker
## -----------------------------
install: build up

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

build:
	$(COMPOSE) build

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

app:
	$(COMPOSE) exec app sh

swagger:
	$(COMPOSE) exec -T app sh -c "swag init -g ./cmd/app/main.go -o ./docs"

postgres:
	$(COMPOSE) exec postgres psql -U $$DB_USER -d $$DB_NAME

create_migration:
	$(COMPOSE) exec app sh -c "migrate create -ext sql -dir /app/migrations -seq ${name}"

migrations_up:
	$(COMPOSE) exec app sh -c "migrate -database $$DB_URL -path /app/migrations up"

clean:
	$(COMPOSE) down -v

prune:
	docker system prune -af --volumes
