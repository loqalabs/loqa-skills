.PHONY: build test lint lint-fix fmt vet clean install-tools help quality-check

# Go parameters
GOCMD=go

# Quality targets for all skills
lint: ## Run golangci-lint on all skills
	@echo "Running golangci-lint on all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Linting $$skill"; \
			cd "$$skill" && golangci-lint run --config ../.golangci.yml; \
			cd ..; \
		fi \
	done

lint-fix: ## Run golangci-lint with auto-fix on all skills
	@echo "Running golangci-lint --fix on all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Fixing $$skill"; \
			cd "$$skill" && golangci-lint run --config ../.golangci.yml --fix; \
			cd ..; \
		fi \
	done

fmt: ## Format code with gofmt on all skills
	@echo "Formatting all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Formatting $$skill"; \
			cd "$$skill" && gofmt -s -w .; \
			cd ..; \
		fi \
	done

vet: ## Run go vet on all skills
	@echo "Running go vet on all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Vetting $$skill"; \
			cd "$$skill" && $(GOCMD) vet ./...; \
			cd ..; \
		fi \
	done

test: ## Run tests on all skills
	@echo "Running tests on all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Testing $$skill"; \
			cd "$$skill" && $(GOCMD) test -v ./...; \
			cd ..; \
		fi \
	done

test-coverage: ## Run tests with coverage on all skills
	@echo "Running tests with coverage on all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Testing $$skill with coverage"; \
			cd "$$skill" && $(GOCMD) test -v -coverprofile=coverage.out ./...; \
			cd ..; \
		fi \
	done

# Development helpers
tidy: ## Tidy go modules for all skills
	@echo "Tidying modules for all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Tidying $$skill"; \
			cd "$$skill" && $(GOCMD) mod tidy; \
			cd ..; \
		fi \
	done

clean: ## Clean build artifacts for all skills
	@echo "Cleaning all Go skills..."
	@for skill in */; do \
		if [ -f "$$skill/go.mod" ]; then \
			echo "Cleaning $$skill"; \
			cd "$$skill" && $(GOCMD) clean && rm -f coverage.out; \
			cd ..; \
		fi \
	done

# Install development tools
install-tools: ## Install development tools
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Pre-commit checks (run before committing)
pre-commit: fmt vet test ## Run all pre-commit checks

# Complete quality validation (run before pushing)
quality-check: fmt vet test lint ## Run comprehensive quality checks on all skills

# Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)