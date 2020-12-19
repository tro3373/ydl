CONTAINERNAME=app
VERSION=1.0.0

all_container=$$(docker ps -a -q)
active_container=$$(docker ps -q)
images=$$(docker images | awk '/^<none>/ { print $$3 }')
local_ip=$$(ip route |awk 'END {print $$NF}')

#.PHONY: all test clean
#default: build
.DEFAULT_GOAL := build

build:
	docker-compose build

up: start
start:
	docker-compose up -d && docker-compose logs -f
stop: down
down:
	docker-compose down
restart: stop start

console: attach
attach:
	docker exec -it $(CONTAINERNAME) /bin/ash --login
do:
	docker exec -it $(CONTAINERNAME) /works/app

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
