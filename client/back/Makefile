.DEFAULT_GOAL := local
local_ip=$$(ip route |grep kernel |awk 'END {print $$NF}')

init:
	@if [[ ! -e node_modules ]]; then \
		echo "======================================= "; \
		echo "==> Executing npm install .."; \
		echo "======================================= "; \
		npm i; \
	fi;

local: init
	@EXTERNAL_IP=$(local_ip) npm run serve

build: init
	@npm run build

# dev: init
# 	@npm run build-dev
#
# prd: init
# 	@npm run build
#
# deploy_dev:
# 	@./bin/deploy
# deploy_prd:
# 	@./bin/deploy -p
