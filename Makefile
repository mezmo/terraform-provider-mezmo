default: build

.PHONY: build
build:
	go build ./...

.PHONY: generate-doc
generate-doc:
	go generate

.PHONY: test
test:
	docker-compose -f compose/base.yml -f compose/test.yml up --remove-orphans --exit-code-from terraform-provider-mezmo --build

.PHONY: local-test
include ./env/local.env
local-test:
	go test -v ./...
