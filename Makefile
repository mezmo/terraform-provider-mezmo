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
	docker-compose -f compose/base.yml pull && docker-compose -p terraform-provider-mezmo-$(BUILD_TAG) -f compose/base.yml -f compose/test.yml up --remove-orphans --exit-code-from terraform-provider-mezmo --build

.PHONY: build-test-example-binary
build-test-example-binary:
	docker run \
	-e GOOS=linux \
	-e GOARCH=amd64 \
	-v $(PWD):/opt/app \
	-w /opt/app \
	goreleaser/goreleaser \
		build \
		--single-target \
		--snapshot \
		--clean \
		-o terraform-provider-mezmo

examples: $(wildcard examples/*/*.tf)

examples/%.tf: build-test-example-binary
	docker run \
		-v $(PWD):/opt/app \
		-v $(PWD)/$@:/opt/example/$(notdir $@) \
		-e TF_CLI_CONFIG_FILE=/opt/app/tf-dev-config \
		-w /opt/example/ \
		hashicorp/terraform:latest \
		validate

.PHONY: local-test
ENV := $(PWD)/env/local.env
include $(ENV)
export
local-test:
	go test -v -count=1 ./...

.PHONY: lint
lint:
	golangci-lint run
