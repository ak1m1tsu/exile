include .env

up:
	@docker compose up -d --build

down:
	@docker compose down --rmi all

gen:
	@go generate ./...

tests:
	@go test -v --race -coverprofile=c.out ./... \
	&& go tool cover -html=c.out \
	&& rm c.out

lint:
	@golangci-lint run ./...
