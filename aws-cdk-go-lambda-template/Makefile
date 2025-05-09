# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

INFRA_DIR := infrastructure
SERVICE_DIR := services/demo
BUILD_LAMBDA_DIR := .build/lambdas

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tools
lint: ## Run linters
	@echo "Running linters..."
	@cd $(SERVICE_DIR) && golangci-lint run --timeout 5m
	@cd $(INFRA_DIR) && golangci-lint run --timeout 5m
.PHONY: lint

test: ## Run tests
	@echo "Running tests..."
	@cd $(SERVICE_DIR) && set -e; \
		go test -v ./... -race -coverprofile=coverage.out && \
		go tool cover -func=coverage.out && \
		rm -f coverage.out
.PHONY: test

##@ Lambda logs
logs:  ### Get logs from the lambda function
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: Lambda name required. Usage: make logs <lambda-name>"; \
		exit 1; \
	fi
	aws logs filter-log-events \
		--log-group-name "/aws/lambda/$(filter-out $@,$(MAKECMDGOALS))" \
		--start-time "$$(date -v -1d +%s000)" \
		--query "events[*].{timestamp:timestamp,message:message}" \
		--output json | \
	jq 'sort_by(.timestamp) | \
		map({ \
			datetime: (.timestamp / 1000 | strftime("%Y-%m-%d %H:%M:%S")), \
			content: (if (.message | startswith("{")) then \
				try (.message | fromjson) catch .message \
				else .message \
				end) \
		})'
.PHONY: logs

# Handle arguments passed to logs target
%:
	@:

##@ Build

build: ## Build the service
	cd $(SERVICE_DIR) && GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o ../../$(BUILD_LAMBDA_DIR)/demo/bootstrap cmd/main.go
.PHONY: build

deploy:  ## Deploy the infrastructure
	cd $(INFRA_DIR) && npx cdk deploy --require-approval=never
.PHONY: deploy

destroy: ## Destroy the infrastructure
	cd $(INFRA_DIR) && npx cdk destroy --force
.PHONY: destroy
