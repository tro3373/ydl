CONTAINER_ydl=ydl
CONTAINER_ngx=nginx
STAGE=dev
arg=

#.PHONY: all test clean
.DEFAULT_GOAL := up

build-ydl:
	@cd ./server/ydl && $(MAKE) build
build_if_needed:
	([[ ! -e ./server/ydl/ydl ]] || $(MAKE) build-ydl); \
		$(MAKE) build-ydl
build-client:
	@cd ./client/back && $(MAKE) build
build: build-ydl build-client

build-image: build
	@docker-compose -f docker-compose.$(STAGE).yml build $(arg)

up: start logsf
start: build_if_needed
	docker-compose -f docker-compose.$(STAGE).yml up -d $(arg)
stop: down
down:
	docker-compose -f docker-compose.$(STAGE).yml down $(arg)
restart:
	docker-compose -f docker-compose.$(STAGE).yml restart $(arg)

logs:
	docker-compose -f docker-compose.$(STAGE).yml logs
logsf:
	docker-compose -f docker-compose.$(STAGE).yml logs -f

console:
	docker exec -it $(CONTAINER_ydl) /bin/sh --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload:
	docker exec -it $(CONTAINER_ngx) nginx -s reload
