CONTAINER_app=ydl
CONTAINER_ngx=nginx
STAGE=dev
arg=

#.PHONY: all test clean
.DEFAULT_GOAL := up

.PHONY: check_ydl
check_ydl:
	@cd ./server/ydl && $(MAKE) check

depends_cmds := docker docker-compose
.PHONY: check
check: check_ydl
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done

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
start: check build_if_needed
	docker-compose -f docker-compose.$(STAGE).yml up -d $(arg)
stop: down
down:
	docker-compose -f docker-compose.$(STAGE).yml down $(arg)
restart:
	docker-compose -f docker-compose.$(STAGE).yml restart $(arg)

logs:
	docker-compose -f docker-compose.$(STAGE).yml logs $(arg)
logsf:
	docker-compose -f docker-compose.$(STAGE).yml logs -f $(arg)

console:
	docker exec -it $(CONTAINER_app)-$(STAGE) /bin/sh --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload-nginx:
	# docker kill $(CONTAINER_ngx)
	docker exec -it $(CONTAINER_ngx) nginx -s reload
