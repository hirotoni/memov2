.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run tests
	go test ./... -cover

.PHONY: depgraph
depgraph: ## Run depgraph
	@echo "Running depgraph"
	goda graph ./... | dot -Tsvg -o depgraph.svg
