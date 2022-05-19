CONTAINER_api=api
CONTAINER_batch=batch
CONTAINER_ngx=nginx
# CONTAINER_db1=mongo001
# CONTAINER_db2=mongo002
# CONTAINER_db3=mongo003

#.PHONY: all test clean
#default: build
.DEFAULT_GOAL := up

build:
	@cd ./server/ydl && $(MAKE) build
build-image: build
	@docker-compose build

up: start logsf
start:
	if [[ ! -e ./server/ydl/ydl ]]; then \
		$(MAKE) build; \
	fi && docker-compose up -d
stop: down
down:
	docker-compose down
restart: stop start
logs:
	docker-compose logs
logsf:
	docker-compose logs -f

console:
	docker exec -it $(CONTAINER_api) /bin/sh --login
console_batch:
	docker exec -it $(CONTAINER_batch) /bin/sh --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
# console_db1:
# 	docker exec -it $(CONTAINER_db1) /bin/bash --login
# console_db2:
# 	docker exec -it $(CONTAINER_db2) /bin/bash --login
# console_db3:
# 	docker exec -it $(CONTAINER_db3) /bin/bash --login
# rs:
# 	docker exec -it $(CONTAINER_db1) /setup_rs
# do:
# 	docker exec -it $(CONTAINER_api) /go/src
