.PHONY: test test-clean

start:
	docker compose up --build

start-daemon:
	docker compose up --build -d

stop:
	docker compose stop

clean:
	docker compose down --remove-orphans --volumes

swagger:
	swag init -g cmd/server/server.go -o api/swagger


test:
	@echo "Running tests..."
	docker compose -f docker-compose.test.yaml up --build --abort-on-container-exit test-runner
	docker compose -f docker-compose.test.yaml down -v

test-clean:
	@echo "Cleaning up test environment..."
	docker compose -f docker-compose.test.yaml down -v
