CURRENT_DIR=$(shell pwd)

swag-init:
	swag init -g api/api.go -o api/docs

start:
	go run cmd/main.go

local-up:
	docker compose --env-file ./.env.docker up -d

local-down:
	docker compose down

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}

pull-sub-module:
	git submodule update --init --recursive

update-sub-module:
	git submodule update --remote --merge

generate:
	go generate ./...

.PHONY: start