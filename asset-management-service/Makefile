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


swag-asset: ##  swag asset-management-service init
	go install github.com/swaggo/swag/cmd/swag@latest
	(swag init -g internal/controller/http/v1/router.go -o docs)
.PHONY: swag-asset


prepare: ##  prepare the environment
	go clean -cache -modcache -i -r
	go clean -modcache && go mod tidy && go mod download
.PHONY: prepare

compose-up: ##  Run docker-compose 
	docker-compose up --build -d  && docker-compose logs -f
.PHONY: compose-up

compose-up-app: ##  ddd
	docker-compose up --build -d asset-management-service && docker-compose logs -f
.PHONY: compose-up-app

compose-up-without-app: ##  Run docker-compose without app
	docker-compose up --build -d wallet-db && docker-compose logs -f
.PHONY: compose-up-without-app



compose-down: ##  Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

run: swag-v1 ## swag run and migrate
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

docker-rm-volume: ##  remove docker volume
	docker volume rm go-cqrs-event-sourcing-tt_pg-data
.PHONY: docker-rm-volume

migrate-create:  ##  create new migration
	migrate create -ext sql -dir migrations 'migrate_name_messages'
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


mock-messages:
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " Legends never die Song by League of legends, Against The Current https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " Canta Per me by Yoki kajiura  https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c","recipient_phone": "05057136988" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " king Song by Florence and the machine https://open.spotify.com/track/1VSngtLdJhrlfHkLxTyOXK?si=d9292df2504e4da0", "recipient_phone": "05057136986" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " mother nature Song by The Hu and LP https://open.spotify.com/track/35SoEGEXsaNnfi8PsT8xEC?si=4347af98187349fa", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " winding river by yu-peng chen https://open.spotify.com/track/04WnFdVesT0VLu1Fc57VoI?si=c12404e1bbb34564","recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " Legends never die Song by League of legends, Against The Current https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " Canta Per me by Yoki kajiura  https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c","recipient_phone": "05057136988" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " king Song by Florence and the machine https://open.spotify.com/track/1VSngtLdJhrlfHkLxTyOXK?si=d9292df2504e4da0", "recipient_phone": "05057136986" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " mother nature Song by The Hu and LP https://open.spotify.com/track/35SoEGEXsaNnfi8PsT8xEC?si=4347af98187349fa", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " winding river by yu-peng chen https://open.spotify.com/track/04WnFdVesT0VLu1Fc57VoI?si=c12404e1bbb34564","recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " Legends never die Song by League of legends, Against The Current https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " Canta Per me by Yoki kajiura  https://open.spotify.com/track/0TI3HDmlvuD0rCwHe5m2wD?si=c8ff85eca986442c","recipient_phone": "05057136988" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " king Song by Florence and the machine https://open.spotify.com/track/1VSngtLdJhrlfHkLxTyOXK?si=d9292df2504e4da0", "recipient_phone": "05057136986" }'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{ "content": " mother nature Song by The Hu and LP https://open.spotify.com/track/35SoEGEXsaNnfi8PsT8xEC?si=4347af98187349fa", "recipient_phone": "05057136986"}'
	curl -X 'POST' 'http://localhost:8080/v1/messages' -H 'accept: application/json' -H 'Content-Type: application/json'  -d '{"content": " winding river by yu-peng chen https://open.spotify.com/track/04WnFdVesT0VLu1Fc57VoI?si=c12404e1bbb34564","recipient_phone": "05057136986"}'

.PHONY mock-messages:


bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

