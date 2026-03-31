# Run application in Docker
.PHONY: docker-run
docker-run:
	docker compose up --force-recreate --build

# Run application locally
.PHONY: local-run
local-run:
	go run ./cmd --config=./config/config.yaml