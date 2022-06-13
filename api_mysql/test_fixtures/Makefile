.PHONY: test

up:
	@echo "start test database"
	docker compose -f docker-compose.test.yml -p mysql_test up -d

stop:
	@echo "stop test database"
	docker compose -f docker-compose.test.yml -p mysql_test stop

down:
	@echo "down test database"
	docker compose -f docker-compose.test.yml -p mysql_test down

test:
	@echo "run test"
	go test ../...

refresh:
	@echo "clean go testcache & recreate db"
	go clean -testcache
	mysql -u $(MYSQL_USER) -h $(MYSQL_HOST) --port=$(MYSQL_PORT) -e "DROP DATABASE $(MYSQL_DATABASE);"
	mysql -u $(MYSQL_USER) -h $(MYSQL_HOST) --port=$(MYSQL_PORT) < ../mock/mysql/docker-entrypoint-initdb.d/1_ddl.sql
	mysql -u $(MYSQL_USER) -h $(MYSQL_HOST) --port=$(MYSQL_PORT) $(MYSQL_DATABASE) < ../mock/mysql/docker-entrypoint-initdb.d/2_data.sql
