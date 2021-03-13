# This file is the development Makefile for the project in this repository.
# All variables listed here are used as substitution in these Makefile targets.

SERVICE_NAME = kafka
AVRO_GENERATOR_PATH = tools/avro-random-generator/

.PHONY: help
help: ## Displays the Makefile help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

################################################################################

RANDOM := $(shell /bin/bash -c "echo $$RANDOM")

.PHONY: generate-mock ## Generate mock AVRO ingestion file to the output folder
generate-mock:
	schema=pkg/avro/transaction_with_amount.avsc
	out=tmp/ingest/payload_$$RANDOM.json
	@if [ -f $(AVRO_GENERATOR_PATH)bin/arg.jar ]; \
	then \
		$(AVRO_GENERATOR_PATH)arg -f pkg/avro/transaction_with_amount.avsc -o tmp/ingest/payload_$$RANDOM.json; \
		echo "Generated"; \
	else \
		echo "Generator not available, building.."; \
		cd $(AVRO_GENERATOR_PATH); \
		./gradlew standalone; \
		echo "Running generator"; \
		./arg -f ../../$(schema) -o ../../$(out); \
		cd -; \
	fi

.PHONY: setup
setup: ## Downloads and install various libs for development.
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	go get golang.org/x/tools/cmd/goimports

.PHONY: mocks ## Generate mocks for interfaces which have been annotated with go generate
mocks:
	go generate ./...

.PHONY: test
test: lint ## Runs the test suite. Some projects might rely on a local development infrasructure to run tests. See `infra-up`.
	go test -v -race -bench=./... -benchmem -timeout=120s -cover -coverprofile=./test/coverage.txt ./...

.PHONY: test-ci
test-ci: lint ## Runs the test suite without -race because it is not supported on alpine
	# This should also run -race but it doesn't work on alpine
	go test -v -benchmem -timeout=120s -cover -coverprofile=./test/coverage.txt -bench=./... ./...

.PHONY: test-quick ## Run the quick test suite
test-quick:
	go test -short -failfast

.PHONY: run
run: build ## Builds project binary and executes it.
	bin/$(SERVICE_NAME)

.PHONY: full
full: clean build fmt lint test ## Cleans up, builds the service, reformats code, lints and runs the test suite.


.PHONY: lint
lint: ## Runs linter against the service codebase.
	golangci-lint run --config configs/golangci.yml
	@echo "[✔] Linter passed"


.PHONY: fmt
fmt: ## Runs gofmt against the service codebase.
	gofmt -w -s .
	goimports -w .
	go clean ./...

.PHONY: tidy
tidy: ## Runs go mod tidy against the service codebase.
	go mod tidy

.PHONY: clean
clean: ## Removes temporary files and deletes the service binary.
	go clean ./...
	rm -f bin/$(SERVICE_NAME)

.PHONY: env
env: ## Displays the current Go environment variables.
	go env

.PHONY: doc
doc: ## Launches the godoc server locally. You can access it through a browser at the URL `http://localhost:8080/`.
	godoc -http=:8080 -index

