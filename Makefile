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
test: test-docs test-unit test-acceptance

.PHONY: test-unit
test-unit:
	go test ./...

.PHONY: test-acceptance
test-acceptance:
	docker-compose -f compose/base.yml pull && docker-compose \
	-p terraform-provider-mezmo-$(BUILD_TAG) \
	-f compose/base.yml -f compose/test.yml \
	up --remove-orphans --exit-code-from terraform-provider-mezmo --build

.PHONY: test-docs
test-docs: generate-doc
	git diff --exit-code -- docs/

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
	goreleaser/goreleaser:v2.0.0 \
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

examples: $(wildcard examples/*/*/)

examples/%: build-test-example-binary
	@echo $@
	@docker run \
		-v $(PWD):/opt/app \
		-v $(PWD)/$@:/opt/example/$@ \
		-e TF_CLI_CONFIG_FILE=/opt/app/tf-dev-config \
		-w /opt/example/$(dir $@) \
		--entrypoint "" \
		hashicorp/terraform:latest \
		/bin/sh -c "[ -d modules ] && terraform init ; terraform validate -compact-warnings" | \
		awk '/Provider development overrides are in effect/ {next} {print}' ; \
		exit $${PIPESTATUS[0]}

.PHONY: local-test
local-test:
	@set -a; . $(PWD)/env/local.env; set +a; \
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run
