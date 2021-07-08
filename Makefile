CONTAINER_api=api
CONTAINER_ngx=nginx
CONTAINER_db1=mongo001
CONTAINER_db2=mongo002
CONTAINER_db3=mongo003

all_container=$$(docker ps -a -q)
active_container=$$(docker ps -q)
images=$$(docker images | awk '/^<none>/ { print $$3 }')
local_ip=$$(ip route |awk 'END {print $$NF}')

#.PHONY: all test clean
#default: build
.DEFAULT_GOAL := up

build:
	docker-compose build

up: start
start:
	docker-compose up -d && docker-compose logs -f
stop: down
down:
	docker-compose down
restart: stop start

console:
	docker exec -it $(CONTAINER_api) /bin/sh --login
console_nginx:
	docker exec -it $(CONTAINER_ngx) /bin/bash --login
console_db1:
	docker exec -it $(CONTAINER_db1) /bin/bash --login
console_db2:
	docker exec -it $(CONTAINER_db2) /bin/bash --login
console_db3:
	docker exec -it $(CONTAINER_db3) /bin/bash --login
# do:
# 	docker exec -it $(CONTAINER_api) /go/src

rs:
	docker exec -it $(CONTAINER_db1) /setup_rs

logs:
	docker-compose logs
logsf:
	docker-compose logs -f

clean: clean_container clean_images
clean_images:
	@if [ "$(images)" != "" ] ; then \
		docker rmi $(images); \
	fi
clean_container:
	@for a in $(all_container) ; do \
		for b in $(active_container) ; do \
			if [ "$${a}" = "$${b}" ] ; then \
				continue 2; \
			fi; \
		done; \
		docker rm $${a}; \
	done
