.PHONY: test

up:
	@echo "start test database"
	docker compose -f docker-compose.test.yml -p mongo_test up -d

stop:
	@echo "stop test database"
	docker compose -f docker-compose.test.yml -p mongo_test stop

down:
	@echo "down test database"
	docker compose -f docker-compose.test.yml -p mongo_test down

test:
	@echo "run test"
	go test ../...