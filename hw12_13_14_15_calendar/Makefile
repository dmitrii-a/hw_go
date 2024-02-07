BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

MIGRATION_DIR="./migrations"
POSTGRES_USER="admin"
POSTGRES_PASSWORD="password"
POSTGRES_DB="calendar-service"
POSTGRES_PORT="5455"
POSTGRES_HOST="127.0.0.1"
POSTGRES_CONTAINER_NAME="calendar-service-db"
GOOSE_DB_STRING="host=${POSTGRES_HOST} port=${POSTGRES_PORT} dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

start-postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) -p $(POSTGRES_PORT):5432 -e POSTGRES_DB=$(POSTGRES_DB) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres:15.4

stop-postgres:
	docker stop $(POSTGRES_CONTAINER_NAME)

rm-postgres:
	docker rm $(POSTGRES_CONTAINER_NAME)

migrate-status:
	goose -dir ${MIGRATION_DIR} postgres ${GOOSE_DB_STRING} status -v

migrate-up:
	goose -dir ${MIGRATION_DIR} postgres ${GOOSE_DB_STRING} up -v

migrate-down:
	goose -dir ${MIGRATION_DIR} postgres ${GOOSE_DB_STRING} down -v

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

lint: install-lint-deps
	golangci-lint run ./...

fmt:
	go fmt ./...
	gofumpt -w -l -extra .
	golangci-lint run ./... --fix
	golines -w .

generate:
	protoc ./api/proto/api/v1/EventService.proto \
			--proto_path=./api/proto \
			--go_out=./internal/presentation/grpc --go_opt=paths=source_relative \
			--go-grpc_out=./internal/presentation/grpc --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=./internal/presentation/grpc --grpc-gateway_opt=paths=source_relative,logtostderr=true \
			--validate_out lang=go:./internal/presentation/grpc --validate_opt=paths=source_relative \
			--swagger_out=logtostderr=true,allow_merge=true,merge_file_name=api:./api \

install-mockery:
	go install github.com/vektra/mockery/v2@v2.40.1

generate-mocks:
	mockery --output=./tests/mocks --all

.PHONY: build run build-img run-img version test lint fix-code-style migrate-up migrate-down migrate-status start-postgres stop-postgres rm-postgres install-lint-deps
