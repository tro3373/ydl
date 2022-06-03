CONTAINER_app=ydl
CONTAINER_ngx=nginx
CONTAINER_client=client
STAGE=dev
arg=

in_app := cd ./server/ydl
in_client := cd ./client/back

.DEFAULT_GOAL := up

check-app:
	@${in_app} && make check
check-client:
	@${in_client} && make check
depends_cmds := docker docker-compose
check: check-app check-client
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done

clean-app:
	@${in_app} && make clean
clean-client:
	@${in_client} && make clean
clean: clean-app clean-client


build-image:
	@docker-compose -f docker-compose.$(STAGE).yml build $(arg)

build-app:
	@docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		app make build
build-client:
	@docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		client make build STAGE=$(STAGE)
build: build-app build-client

up: start logsf
start:
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
console_client:
	docker exec -it $(CONTAINER_client) /bin/bash --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload-nginx:
	# docker kill $(CONTAINER_ngx)
	docker exec -it $(CONTAINER_ngx) nginx -s reload
