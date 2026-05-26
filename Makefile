DOCKER_BUILD = docker build -f

.PHONY: docker-build-auth docker-build-gateway docker-build-user
.PHONY: docker-build-comment docker-build-file docker-build-interaction
.PHONY: docker-build-post docker-build-search docker-build-section
.PHONY: docker-build-all

docker-build-auth:
	$(DOCKER_BUILD) auth/Dockerfile -t temp-auth .

docker-build-gateway:
	$(DOCKER_BUILD) gateway/Dockerfile -t temp-gateway .

docker-build-user:
	$(DOCKER_BUILD) user/Dockerfile -t temp-user .

docker-build-comment:
	$(DOCKER_BUILD) comment/rpc/Dockerfile -t temp-comment .

docker-build-file:
	$(DOCKER_BUILD) file/rpc/Dockerfile -t temp-file .

docker-build-interaction:
	$(DOCKER_BUILD) interaction/rpc/Dockerfile -t temp-interaction .

docker-build-post:
	$(DOCKER_BUILD) post/rpc/Dockerfile -t temp-post .

docker-build-search:
	$(DOCKER_BUILD) search/rpc/Dockerfile -t temp-search .

docker-build-section:
	$(DOCKER_BUILD) section/rpc/Dockerfile -t temp-section .

docker-build-all: docker-build-auth docker-build-gateway docker-build-user \
                  docker-build-comment docker-build-file \
                  docker-build-interaction docker-build-post \
                  docker-build-search docker-build-section
