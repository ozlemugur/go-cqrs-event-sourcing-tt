include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


swag-wallet: ## swag wallet-management-service init 
	go install github.com/swaggo/swag/cmd/swag@latest
	( swag init -g internal/controller/http/v1/router.go -o docs )
.PHONY: swag-wallet

prepare: ##  prepare the environment
	go clean -cache -modcache -i -r
	go clean -modcache && go mod tidy && go mod download && go mod vendor
.PHONY: prepare

compose-up: ##  Run docker-compose 
	docker-compose up --build -d  && docker-compose logs -f
.PHONY: compose-up

compose-up-app: ##  ddd
	docker-compose up --build -d wallet-management-service && docker-compose logs -f
.PHONY: compose-up-app

compose-down: ##  Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down


migrate-create: ##  create new migration from sql
	migrate create -ext sql -dir migrations 'wallet_management_service'
.PHONY: migrate-create

migrate-up: ##  migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up



integration-test: ##  run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

check-health:## check health
	curl -f http://localhost:8080/healthz || echo "Health check failed"
.PHONY:check-health


bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

