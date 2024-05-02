BUILD_TAG ?= "1"
SHELL = '/bin/bash'

default: build

.PHONY: build
build:
	go build ./...

.PHONY: generate-doc
generate-doc:
	go generate

.PHONY: test
test:
	docker-compose -f compose/base.yml pull && docker-compose \
	-p terraform-provider-mezmo-$(BUILD_TAG) \
	-f compose/base.yml -f compose/test.yml \
	up --remove-orphans --exit-code-from terraform-provider-mezmo --build

.PHONY: start
start:
	docker-compose -f compose/base.yml pull && docker-compose \
	-p mezmo-tf \
	-f compose/base.yml -f compose/dev.yml \
	up --remove-orphans --build

# Builds a binary using goreleaser. Uses the systems GOOS/GOARCH unless one is
# already specified.
.PHONY: build-snapshot-binary
build-snapshot-binary:
	$(eval GOOS ?= $(shell go env GOOS))
	$(eval GOARCH ?= $(shell go env GOARCH))
	@docker run \
	-e GOOS \
	-e GOARCH \
	-v $(PWD):/opt/app \
	-w /opt/app \
	goreleaser/goreleaser \
		build \
		--single-target \
		--snapshot \
		--clean \
		-o terraform-provider-mezmo

# Set vars for the test-example target to be build
.PHONY: set-test-example-target
set-test-example-target:
	$(eval GOOS=linux)
	$(eval GOARCH=amd64)
	$(eval export)

# Build a binary to be used to test the examples
.PHONY: build-test-example-binary
build-test-example-binary: set-test-example-target build-snapshot-binary

examples: $(wildcard examples/*/*/*.tf)

examples/%.tf: build-test-example-binary
	@docker run \
		-v $(PWD):/opt/app \
		-v $(PWD)/$@:/opt/example/$(notdir $@) \
		-e TF_CLI_CONFIG_FILE=/opt/app/tf-dev-config \
		-w /opt/example/ \
		hashicorp/terraform:latest \
		validate -compact-warnings | \
		awk '/Provider development overrides are in effect/ {next} {print}' ; \
		exit $${PIPESTATUS[0]}

.PHONY: local-test
ENV := $(PWD)/env/local.env
include $(ENV)
export
local-test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run
