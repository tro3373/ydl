CONTAINER_ngx := nginx
CONTAINER_client := client
export OWNER := $(if $(OWNER),$(OWNER),$(shell git config --get remote.origin.url |sed -e 's,^.*:,,g' -e 's,/.*,,g'))
export APP_NAME := ydl
export APP_VER := 1.0.0

DCY := docker-compose.yml
export STAGE := $(if $(STAGE),$(STAGE),$(shell if [[ -e $(DCY) ]] && readlink $(DCY) |grep prd.yml>&/dev/null; then echo prd; else echo dev; fi))
ARG :=

IN_APP := cd ./server/ydl
IN_CLIENT := cd ./client/back

.DEFAULT_GOAL := up

tag:
	@tag="v${APP_VER}" && git tag "$$tag" && echo "==> $$tag tagged."

check-app:
	@${IN_APP} && make check
check-client:
	@${IN_CLIENT} && make check
depends_cmds := docker docker-compose
check:
	@test -e $(DCY) || (echo "==> No $(DCY) exist" 1>&2 && exit 1)
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done
check-all: check-app check-client check

clean-app:
	@${IN_APP} && make clean
clean-client:
	@${IN_CLIENT} && make clean
clean: clean-app clean-client


build-image-dev: _build-image-dev
build-image-prd: _build-image-prd
_build-image-%:
	@echo "==> $@"
	@docker-compose -f docker-compose.${*}.yml build $(ARG)
	@echo "Done"

push-image:
	@docker push $(OWNER)/$(APP_NAME):$(APP_VER)
pull-image:
	@docker pull $(OWNER)/$(APP_NAME):$(APP_VER)

build-client:
	@docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		client make build STAGE=prd
build-app:
	@docker-compose -f docker-compose.dev.yml \
		run --rm -it \
		app make build
build: build-client build-app

prepare: check
	@if [[ ${STAGE} == "dev" ]]; then \
		docker images |grep ydl-dev >&/dev/null || make build-image-dev; \
	fi

up: start logsf
start: prepare
	@docker-compose -f $(DCY) up -d $(ARG)
stop: down
down:
	@docker-compose -f $(DCY) down $(ARG)
restart:
	@docker-compose -f $(DCY) restart $(ARG)

logs:
	@docker-compose -f $(DCY) logs $(ARG)
logsf:
	@docker-compose -f $(DCY) logs -f $(ARG)

console:
	@docker exec -it $(APP_NAME)-app /bin/sh --login
console_client:
	@docker exec -it $(CONTAINER_client) /bin/bash --login
console_nginx:
	@docker exec -it $(CONTAINER_ngx) /bin/bash --login
reload-nginx:
	@docker exec -it $(CONTAINER_ngx) nginx -s reload

link-dev: _link-dev
link-prd: _link-prd
_link-%:
	@if [[ -e $(DCY) ]]; then rm $(DCY); fi
	@ln -s docker-compose.${*}.yml $(DCY)
	@echo "Done. linked to $$(readlink $(DCY))"
