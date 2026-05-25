DOCKER_BUILD = docker build --build-arg SERVICE

.PHONY: docker-build-auth docker-build-gateway docker-build-user
.PHONY: docker-build-comment docker-build-file docker-build-interaction
.PHONY: docker-build-post docker-build-search docker-build-section
.PHONY: docker-build-all

docker-build-auth:
	$(DOCKER_BUILD)=auth -t temp-auth .

docker-build-gateway:
	$(DOCKER_BUILD)=gateway -t temp-gateway .

docker-build-user:
	$(DOCKER_BUILD)=user -t temp-user .

docker-build-comment:
	$(DOCKER_BUILD)=comment/rpc -t temp-comment .

docker-build-file:
	$(DOCKER_BUILD)=file/rpc -t temp-file .

docker-build-interaction:
	$(DOCKER_BUILD)=interaction/rpc -t temp-interaction .

docker-build-post:
	$(DOCKER_BUILD)=post/rpc -t temp-post .

docker-build-search:
	$(DOCKER_BUILD)=search/rpc -t temp-search .

docker-build-section:
	$(DOCKER_BUILD)=section/rpc -t temp-section .

docker-build-all: docker-build-auth docker-build-gateway docker-build-user \
                  docker-build-comment docker-build-file \
                  docker-build-interaction docker-build-post \
                  docker-build-search docker-build-section
