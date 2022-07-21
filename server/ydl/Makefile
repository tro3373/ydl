SHELL=/bin/bash
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
name := $(shell grep module $(current_dir)/go.mod|head -1|sed -e 's,^.*/,,g')

.DEFAULT_GOAL := run

depends_cmds := go gosec statik
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done
	@echo "[OK] check ok!"

clean:
	@for d in $(name); do if [[ -e $${d} ]]; then echo "==> Removing $${d}.." && rm -rf $${d}; fi done
	@echo "[OK] clean ok!"

run: check clean generate
	@go run ./main.go

generate:
	@statik -f -src=./assets/
	@echo "[OK] Files added to /statik!"

generate-if-needed:
	@if [[ ! -e ./statik/statik.go ]]; then \
		make generate; \
	fi

sec:
	@gosec --color=false ./...
	@echo "[OK] Go security check was completed!"

# build: check clean generate-if-needed sec
# 	@env GOOS=linux go build -ldflags="-s -w"
build:
	@pre_env="env GOOS=linux " make _build

# build-android: check clean generate-if-needed sec
# 	@env GOOS=android GOARCH=arm64 go build -ldflags="-s -w"
# 	@echo "is api url :3000?"
# 	@echo "Do in termux, pkg install proot && termux-chroot && echo 'nameserver 8.8.8.8' >/etc/resolv.conf"
build-android:
	@pre_env="env GOOS=android GOARCH=arm64" make _build

_build: check clean generate-if-needed sec
	@$(pre_env) go build -ldflags="-s -w"

deps:
	@go list -m all

tidy:
	@go mod tidy

curl_get:
	@curl \
		-fSsL \
		-H 'Content-Type: application/json' \
		http://127.0.0.1:3000/api

curl_post:
	@curl \
		-fSsL \
		-X POST \
		-H 'Content-Type: application/json' \
		-d@sample.json \
		http://127.0.0.1:3000/api
