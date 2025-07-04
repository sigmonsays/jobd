.PHONY: help
help: ## Show this help screen
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

gen: ## gen api
	go generate ./api

g: gen ## alias for gen

deps: ## Install deps
	go install github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/cespare/reflex@latest
