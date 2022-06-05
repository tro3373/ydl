SHELL=/bin/bash
STAGE=dev
local_ip=$$(ip route |grep kernel |awk 'END {print $$NF}')

.DEFAULT_GOAL := run

depends_cmds := npm
check:
	@for cmd in ${depends_cmds}; do command -v $$cmd >&/dev/null || (echo "No $$cmd command" && exit 1); done

clean:
	@for d in node_modules dist; do if [[ -e $${d} ]]; then echo "==> Removing $${d}.." && rm -rf $${d}; fi done

init:
	@if [[ ! -e node_modules ]]; then \
		echo "======================================= "; \
		echo "==> Executing npm i.."; \
		echo "======================================= "; \
		npm i; \
	fi;

run: check init
	@EXTERNAL_IP=$(local_ip) npm run serve-$(STAGE)

build: check init
	@npm run build-$(STAGE)
