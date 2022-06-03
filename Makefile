exec:
	docker exec -it know-sync-api sh

build:
	docker-compose build --no-cache
	
up:
	docker-compose up

down:
	docker-compose down

test:
	go test -v ./cmd/...

.PHONY: exec build up down test
