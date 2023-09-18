up:
	./scripts/dockerup.sh

down:
	docker compose down --rmi all

gen:
	go generate ./...

tests:
	go test -v --race -coverprofile=c.out ./... \
	&& go tool cover -html=c.out \
	&& rm c.out

seedkafka:
	./scripts/seedkafka.sh
