BUILD_TAG ?= "1"

default: build

.PHONY: build
build:
	go build ./...

.PHONY: generate-doc
generate-doc:
	go generate

.PHONY: test
test:
	printenv
	docker-compose -p terraform-provider-mezmo-$(BUILD_TAG) -f compose/base.yml -f compose/test.yml up --remove-orphans --exit-code-from terraform-provider-mezmo --build

.PHONY: local-test
ENV := $(PWD)/env/local.env
include $(ENV)
export
local-test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run
