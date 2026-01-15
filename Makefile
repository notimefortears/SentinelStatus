# Sentinel Makefile

.PHONY: start stop restart clean logs

start:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	docker-compose down
	docker-compose up -d --build

clean:
	docker-compose down -v
	docker system prune -f

logs:
	docker-compose logs -f worker

# Shortcut to add a test target
test-target:
	curl -X POST http://localhost:3000/targets \
		-H "Content-Type: application/json" \
		-d '{"url": "https://google.com"}'