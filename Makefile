# Define some colors for help output
NC=\033[0m
BOLD=\033[1m
GREEN=\033[32m

.PHONY: re help

re: ## Build the application
	@echo "Building the application..."
	go build -o re cmd/re/main.go 
	@echo "Build complete."

help: ## Display this help message
	@printf "\n$(BOLD)Available targets:$(NC)\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BOLD)$(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@printf "\n"
