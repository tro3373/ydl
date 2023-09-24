SHELL=/bin/bash
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
name := $(shell grep module $(current_dir)/go.mod|head -1|sed -e 's,^.*/,,g')
dist := ../../client/back/dist
srcf := $(dist)/index.html
dstf := ./statik/statik.go

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
	@if [[ ! -e $(srcf) ]]; then \
		echo "==> Client build needed." 1>&2; \
		exit 1; \
	fi
	@if [[ ! -e $(dstf) || $(srcf) -nt $(dstf) ]]; then \
		echo "==> Generating statik..."; \
		statik -f -src=$(dist); \
		echo "[OK] Files added to /statik!"; \
	fi

sec:
	@gosec --color=false ./...
	@echo "[OK] Go security check was completed!"

build:
	@pre_env="env GOOS=linux " make _build

build-android:
	@pre_env="env GOOS=android GOARCH=arm64" make _build

_build: check clean generate sec
	@$(pre_env) go build -ldflags="-s -w"

deps:
	@go list -m all

tidy:
	@go mod tidy

tidy-go:
	@ver=$(shell go version|awk '{print $$3}' |sed -e 's,go\(.*\)\..*,\1,g') && go mod tidy -go=$${ver}

gr_check:
	@goreleaser check
gr_snap:
	@goreleaser release --snapshot --rm-dist $(OPT)
gr_snap_skip_publish:
	@OPT=--skip-publish make gr_snap
gr_build:
	@goreleaser build --snapshot --rm-dist

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
