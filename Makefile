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
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done
check-all: check-app check-client check

clean-app:
	@${in_app} && make clean
clean-client:
	@${in_client} && make clean
clean: clean-app clean-client


build-image:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.$(STAGE).yml build $(arg)

build-app:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		app make build
build-client:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		client make build STAGE=$(STAGE)
build: build-app build-client

prepare: check
	@echo "==> $@ $(STAGE)" && \
	(docker images |grep ydl-dev >&/dev/null || make STAGE=dev build-image) && \
	if [[ ${STAGE} == "prd" ]]; then \
		([[ ! -e ./client/back/dist || ! -e ./server/ydl/ydl ]] && make STAGE=prd build); \
		(docker images |grep 'ydl ' >&/dev/null || make STAGE=prd build-image) \
	fi

up: start logsf
start: prepare
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
	docker exec -it $(CONTAINER_app)-app /bin/sh --login
console_client:
	docker exec -it $(CONTAINER_client) /bin/bash --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload-nginx:
	# docker kill $(CONTAINER_ngx)
	docker exec -it $(CONTAINER_ngx) nginx -s reload
