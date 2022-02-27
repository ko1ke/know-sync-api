attach:
	docker exec -it know-sync-api bash

build:
	docker-compose build --no-cache
	
up:
	docker-compose up

down:
	docker-compose down

.PHONY: attach up down
