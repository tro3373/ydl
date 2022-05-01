# @see
# https://qiita.com/Syoitu/items/8e7e3215fb7ac9dabc3a
# https://qiita.com/keitakn/items/f46347f871083356149b
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
name := $(shell head -1 $(current_dir)/go.mod|sed -e 's,^.*/,,g')
version := v0.0.1

.DEFAULT_GOAL := run

.PHONY: deps
deps:
	@go list -m all

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: build
build:
	@env GOOS=linux go build -ldflags="-s -w"
.PHONY: build-image
build-image:
	@$(MAKE) -f ../../Makefile build

.PHONY: clean
clean:
	@[[ -e $(name) ]] && rm -rfv $(name)

.PHONY: run
run:
	@go run ./main.go |& zap-pretty
.PHONY: rund
rund:
	docker run --rm -it \
		-v $(current_dir):/go/src \
		-w /go/src \
		-p 3000:3000 \
		go-air:1.0.0

.PHONY: curl_get
curl_get:
	@curl \
		-fSsL \
		-H 'Content-Type: application/json' \
		http://127.0.0.1:3000/api

.PHONY: curl_post
curl_post:
	@curl \
		-fSsL \
		-X POST \
		-H 'Content-Type: application/json' \
		-d@sample.json \
		http://127.0.0.1:3000/api
