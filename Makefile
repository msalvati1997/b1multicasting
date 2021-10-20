start: ## Start docker development environment
	docker-compose up
down: clear ## Destroy the development environment
	docker-compose down
	rm -rf var/docker/volumes/*
stop: ## Stop docker development environment
	docker-compose stop

