CONTAINER_app := ydl
CONTAINER_ngx := nginx
CONTAINER_client := client
STAGE := dev
ARG :=

IN_APP := cd ./server/ydl
IN_CLIENT := cd ./client/back

.DEFAULT_GOAL := up

check-app:
	@${IN_APP} && make check
check-client:
	@${IN_CLIENT} && make check
depends_cmds := docker docker-compose
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done
check-all: check-app check-client check

clean-app:
	@${IN_APP} && make clean
clean-client:
	@${IN_CLIENT} && make clean
clean: clean-app clean-client


build-image:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.$(STAGE).yml build $(ARG)

build-client:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		client make build STAGE=$(STAGE)
build-app:
	@echo "==> $@ $(STAGE)" && \
	docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		app make build
build: build-client build-app

prepare: check
	@echo "==> $@ $(STAGE)" && \
	(docker images |grep ydl-dev >&/dev/null || make STAGE=dev build-image) && \
	if [[ ${STAGE} == "prd" ]]; then \
		(docker images |grep 'ydl ' >&/dev/null || make STAGE=prd build-image) && \
		(test -e ./client/back/dist || make STAGE=prd build-client) && \
		(test -e ./server/ydl/ydl || make STAGE=prd build-app); \
	fi

up: start logsf
start: prepare
	docker-compose -f docker-compose.$(STAGE).yml up -d $(ARG)
stop: down
down:
	docker-compose -f docker-compose.$(STAGE).yml down $(ARG)
restart:
	docker-compose -f docker-compose.$(STAGE).yml restart $(ARG)

logs:
	docker-compose -f docker-compose.$(STAGE).yml logs $(ARG)
logsf:
	docker-compose -f docker-compose.$(STAGE).yml logs -f $(ARG)

console:
	docker exec -it $(CONTAINER_app)-app /bin/sh --login
console_client:
	docker exec -it $(CONTAINER_client) /bin/bash --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload-nginx:
	# docker kill $(CONTAINER_ngx)
	docker exec -it $(CONTAINER_ngx) nginx -s reload
