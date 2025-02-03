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
	( cd wallet-management-service && swag init -g internal/controller/http/v1/router.go -o docs )
.PHONY: swag-wallet

swag-asset: ##  swag asset-management-service init
	go install github.com/swaggo/swag/cmd/swag@latest
	(cd asset-management-service && swag init -g internal/controller/http/v1/router.go -o docs)
.PHONY: swag-asset

swag-query: ##  swag asset-query-service init
	go install github.com/swaggo/swag/cmd/swag@latest
	(cd asset-query-service && swag init -g internal/controller/http/v1/router.go -o docs)
.PHONY: swag-query

prepare: ##  prepare the environment
	go clean -cache -modcache -i -r
	go clean -modcache && go mod tidy && go mod download
.PHONY: prepare

compose-up: ##  Run docker-compose 
	docker-compose up --build -d  wallet-db  query-db kafka wallet-management-service   asset-query-service   asset-query-processor  asset-management-service  asset-processor  kafdrop && docker-compose logs -f
.PHONY: compose-up

compose-up-db: ##  Run docker-compose databases
	docker-compose up --build -d  wallet-db  query-db kafka && docker-compose logs -f
.PHONY: compose-up

compose-up-app: ##  compose up just apss
	docker-compose build --no-cache  wallet-management-service  asset-query-service   asset-query-processor  asset-management-service
	docker-compose up  -d wallet-management-service  asset-query-service   asset-query-processor  asset-management-service   && docker-compose logs -f
.PHONY: compose-up-app



compose-down: ##  Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

remove-volumes: ##  remove docker volume
	docker volume ls | grep go-cqrs-event-sourcing-tt | awk '{print $2}' | xargs docker volume rm
.PHONY: docker-rm-volume

login-wallet-db: ## Login to wallet-db, list tables, and display rows from wallets table
	docker exec -it wallet-db psql -U wallet_user -d wallet_db -c "\dt"
	docker exec -it wallet-db psql -U wallet_user -d wallet_db -c "SELECT * FROM wallets;"
.PHONY: login-wallet-db


login-query-db: ## Login to query-db, list tables, and display rows from wallets table
	docker exec -it query-db psql -U query_user -d query_db -c "\dt"
	docker exec -it query-db psql -U query_user -d query_db -c  "SELECT * FROM wallet_assets;"
.PHONY: login-query-db

# Name of your project images (prefix or full name)
PROJECT_IMAGE_PREFIX := go-cqrs-event-sourcing-tt_

# Removes dangling images and unused images that are not related to the project
clean-docker-images: ## Remove unused Docker images, excluding project images
	@echo "Removing dangling images..."
	docker images -f "dangling=true" -q | xargs -r docker rmi
	@echo "Removing unused images not related to the project..."
	docker images | grep -v "$(PROJECT_IMAGE_PREFIX)" | awk 'NR>1 {print $$3}' | xargs -r docker rmi -f

.PHONY: clean-docker-images



migrate-create:  ##  create new migration from sql
	migrate create -ext sql -dir migrations 'migrate_name_messages'
.PHONY: migrate-create

migrate-up: ##  migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

integration-test: ##  run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

create-wallets: ## Send three POST requests to create wallets
	curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "CrookedStil", "network": "Bitcoin"}'

	curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "CrookedStil1", "network": "Bitcoin"}'

	curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "CrookedStil2", "network": "Bitcoin"}'

	curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "CrookedStil3", "network": "Bitcoin"}'\

	  curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "LP-TheHU-MotherNature", "network": "Ethereum"}'\

	  curl -X 'POST' \
	  'http://localhost:8081/v1/wallets' \
	  -H 'accept: application/json' \
	  -H 'Content-Type: application/json' \
	  -d '{"address": "LP-TheHU-MotherNature-2", "network": "Ethereum"}'\
	  && make login-wallet-db
.PHONY: create-wallets


deposit-wallets: ## deposit
	curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 10, "asset_name": "BTC", "wallet_id": 1}'

	curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 20, "asset_name": "BTC", "wallet_id": 2}'

	curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 30, "asset_name": "BTC", "wallet_id": 3}'
	
	curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 40, "asset_name": "BTC", "wallet_id": 4}'\

		curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 40, "asset_name": "BTC", "wallet_id": 5}'\

		curl -X 'POST' \
	'http://localhost:8082/v1/assets/deposit' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 40, "asset_name": "BTC", "wallet_id": 6}'\
	&& make login-query-db

.PHONY: deposit-wallets 


withdraw-wallet: ## Withdraw amount from a wallet
	make login-query-db  && \
    curl -X 'POST' \
	'http://localhost:8082/v1/assets/withdraw' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{ "amount": 20, "asset_name": "BTC", "wallet_id": 1 }' \
	&& make login-query-db

.PHONY: withdraw-wallet


EPOCH_15_MINUTES_LATER := $(shell date -v+15M +%s)


transfer-wallet: ## Transfer funds with dynamic execute_time for macos solution
	curl -X 'POST' \
	'http://localhost:8082/v1/assets/transfer' \
	-H 'accept: application/json' \
	-H 'Content-Type: application/json' \
	-d '{"amount": 120, "asset_name": "BTC", "execute_time": $(EPOCH_15_MINUTES_LATER), "from_wallet_id": 3, "to_wallet_id": 2}' \
	&& make login-query-db
	
.PHONY: transfer-wallet





bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

