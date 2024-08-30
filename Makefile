export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)



# Полезные скрипты

swag: ### swagger init
	 swag init -g main.go
.PHONY: swag


docker-build: ### build app in docker
	docker build -t test_task_kami:latest -f Dockerfile .
.PHONY: docker-build


compose-up: ### build app in docker
	docker-compose up
.PHONY: compose-up


migration-up: ### migration up
	@source .env.local && \
	DB_URL="postgresql://$${DB_USERNAME}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSL_MODE}" && \
	goose -dir internal/migrations postgres "$$DB_URL" up
.PHONY: migration-up

migration-down: ### migration down
	@source .env.local && \
	DB_URL="postgresql://$${DB_USERNAME}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSL_MODE}" && \
	goose -dir internal/migrations postgres "$$DB_URL" down
.PHONY: migration-down

migration-create:
	goose -dir internal/migrations create $(name) sql
.PHONY: migration-create

test:
	 cd internal/handler && go test -v
.PHONY: test



