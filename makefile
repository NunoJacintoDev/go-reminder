.PHONY: docker
docker:
	docker-compose run --service-ports app bash
