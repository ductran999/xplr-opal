default: help

help: ## Show help for each of the Makefile commands
	@awk 'BEGIN \
		{FS = ":.*##"; printf "Usage: make ${cyan}<command>\n${white}Commands:\n"} \
		/^[a-zA-Z_-]+:.*?##/ \
		{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
		$(MAKEFILE_LIST)


.PHONY: demo
setup: ## Setup keycloak, opal for demo
	@docker compose up -d

.PHONY: run
run: ## Run app 
	go run cmd/main.go

.PHONY: api
api: ## Auto generate api code from openapi.yml
	@echo "==> Generating API code from OpenAPI spec"
	mkdir -p ./api/generated
	rm -f ./api/generated/*
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	go generate ./...
	go mod tidy
