# Mezmo Terraform Provider

The Mezmo Terraform Provider allows organizations to manage Pipelines (sources, processors and destinations)
programmatically via Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

```shell
go build ./...
```

### Running unit tests

```shell
go test ./...
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real pipeline resources against the configured environment

```shell
make testacc
```
