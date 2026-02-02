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
	@echo "  make build            🔨 Build containers"
	@echo "  make restart          🔁 Restart containers"
	@echo "  make logs             📜 Show logs"
	@echo "  make ps               📦 Show containers"
	@echo "  make app              🐹 Enter app container"
	@echo "  make postgres         🐘 Enter postgres container"
	@echo "  make create_migration 📝 Create migrations"
	@echo "  make migrations_up    ⬆️ Run migrations"
	@echo "  make clean            🧹 Remove containers + volumes"
	@echo "  make prune            💣 Docker system prune"
	@echo ""


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

postgres:
	$(COMPOSE) exec postgres psql -U $$DB_USER -d $$DB_NAME

create_migration:
	migrate create -ext sql -dir src/migrations -seq ${name}

migrations_up:
	$(COMPOSE) exec app sh -c "migrate -database $$DB_URL -path src/migrations up"

clean:
	$(COMPOSE) down -v

prune:
	docker system prune -af --volumes
