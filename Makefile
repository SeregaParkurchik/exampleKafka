restart:
	docker compose -f docker-compose/docker-compose-dev.yaml down && \
	docker compose -f docker-compose/docker-compose-dev.yaml up --build -d
	@echo "--- Контейнеры успешно перезапущены! ---"

