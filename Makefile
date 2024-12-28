setup:
	go mod tidy
	go mod vendor
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install github.com/evilmartians/lefthook@latest
	lefthook install
	git config --local core.hooksPath .git/hooks
	go install github.com/rubenv/sql-migrate/...@latest
	go install go.uber.org/mock/mockgen@latest

lint:
	golangci-lint -c ./tools/.golangci.yml run --fix

GO_ENV ?= development
run: migrate-up
	GO_ENV=$(GO_ENV) go run ./cmd/main.go

migrate-up:
	sql-migrate up -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

migrate-down:
	sql-migrate down -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

migrate-status:
	sql-migrate status -config=./tools/db/dbconfig.yml -env=$(GO_ENV)

new-migration:
	@read -p "Migration name: " NAME; \
    	if [ -z "$$NAME" ]; then \
    		echo "Invalid migration name"; \
    		exit 1; \
    	fi; \
    	sql-migrate new -config=./tools/db/dbconfig.yml $$NAME
test:
	go test ./...

mock: mock-ports mock-domain

MOCK_PORTS := \
	wordports/queue \
	wordports/scrapper

mock-ports:
	$(foreach file,$(MOCK_PORTS),\
		mockgen -typed=true \
		-source=./internal/ports/$(file).go \
		-destination=./internal/ports/$(dir $(file))mock/$(notdir $(file)).go;)

MOCK_DOMAIN := \
	word/service

mock-domain:
	$(foreach file,$(MOCK_DOMAIN),\
		mockgen -typed=true \
		-source=./internal/domain/$(file).go \
		-destination=./internal/domain/$(dir $(file))mock/$(notdir $(file)).go;)
