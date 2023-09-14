# Mezmo Terraform Provider

The Mezmo Terraform Provider allows organizations to manage Pipelines (sources, processors and destinations)
programmatically via Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building the Provider

```shell
go build ./...
```

## Running Integration Tests

### Within Docker

```shell
make test
```

### Locally

Start the services:

```shell
docker-compose -f compose/base.yml pull && docker-compose -f compose/base.yml up --remove-orphans
```

Running the test suite:

```shell
make local-test
```

Running a single test

```shell
env $(cat ./env/local.env) go test -v -count=1 ./... -run TestDedupeProcessor
```

## Generating the Docs

To generate or update documentation, run `go generate`.

## Using the Provider

To install the provider in development, run `go install .`. This will build the provider and put the provider
binary in the `$GOPATH/bin` directory.

Then, you can `plan` or `apply` a terraform files:

```bash
pushd examples/pipeline
terraform plan
popd
```

