SQLC_VERSION ?= v1.30.0
BINARY ?= service

.PHONY: help sqlc gqlgen run build test vet lint rename tidy clean compose-up compose-down compose-build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## sqlc: Generate Go code from SQL queries (internal/infrastructure/database/pg/queries)
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

## gqlgen: Regenerate GraphQL code after schema changes
gqlgen:
	go tool gqlgen generate

## run: Run the API in development mode (migrations are applied on startup)
run: sqlc
	go run ./cmd/api

## build: Build the API binary
build: sqlc
	go build -o $(BINARY) ./cmd/api

## test: Run all tests
test: sqlc
	go test ./...

## vet: Run static analysis
vet: sqlc
	go vet ./...

## lint: Run golangci-lint (https://golangci-lint.run/welcome/install/)
lint: sqlc
	golangci-lint run

## rename: Rename the Go module, e.g. make rename MODULE=github.com/your-org/my-service
rename: sqlc
	@test -n "$(MODULE)" || { echo "usage: make rename MODULE=github.com/your-org/my-service"; exit 1; }
	@OLD=$$(go list -m) && \
	grep -rlF "$$OLD" . --include='*.go' --include='go.mod' --include='*.yml' --include='*.yaml' --include='*.md' --include='Dockerfile' --exclude-dir=.git \
		| xargs sed -i.bak "s|$$OLD|$(MODULE)|g" && \
	find . -name '*.bak' -not -path './.git/*' -delete && \
	go mod tidy && echo "Module renamed: $$OLD -> $(MODULE)"

## tidy: Tidy Go module dependencies
tidy:
	go mod tidy

## clean: Remove the binary and generated sqlc code
clean:
	rm -f $(BINARY)
	rm -rf internal/infrastructure/database/pg/sqlc

## compose-up: Start the stack (API + PostgreSQL) with Docker Compose
compose-up:
	docker compose up --build

## compose-down: Stop the Docker Compose stack
compose-down:
	docker compose down
