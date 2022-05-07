-include: .env

GO_ENTRYPOINT=cmd/main.go
COVERAGE_OUTPUT=coverage.output coverage.out
COVERAGE_HTML=coverage.html
GO_PACKAGES=cmd internal pkg

## @ Help
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make [target]\033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)


## @ Application
.PHONY: run setup setdown
run: ## Start application
	@go run $(GO_ENTRYPOINT)

setup: ## Start app dependencies
	@docker-compose up -d

setdown:  ## Stop application dependencies
	@docker-compose down


## @ Linter
.PHONY: lint format
lint:
	@golangci-lint run -v

format:
	@gofumpt -w -e -l $(GO_PACKAGES)

## @ Create docs
docs: cmd internal pkg
	@swag init -g $(GO_ENTRYPOINT)
	@swag fmt -g $(GO_ENTRYPOINT)


## @ Tests
.PHONY: test coverage
test:  ## Run tests of project
	@go test ./... -race -v

coverage: ## Run tests, make report and open into browser
	@go test ./... -race -v -cover  -covermode=atomic -coverprofile=$(COVERAGE_OUTPUT)
	@go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_HTML)
	@wslview ./$(COVERAGE_HTML) || xdg-open ./$(COVERAGE_HTML) || powershell.exe Invoke-Expression ./$(COVERAGE_HTML)

## @ Clean
.PHONY: clean clean_coverage_cache
clean_coverage_cache: ## Remove coverage cache files
	@rm -rf $(COVERAGE_OUTPUT)
	@rm -rf $(COVERAGE_HTML)

clean: clean_coverage_cache ## Remove Cache files
