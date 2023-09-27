# Mezmo Terraform Provider

The Mezmo Terraform Provider allows organizations to manage Pipelines (sources, processors and destinations)
programmatically via Terraform.

You can download this repo to create your own local provider or you can use the Hashicorp registry.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building the Provider Locally

```shell
go build ./...
```

## Generating the Docs

To generate or update documentation, run `go generate`.

## Using the Provider

To install the provider in development, run `go install .`. This will build the provider and put the provider
binary in the `$GOPATH/bin` directory.

## Adding the Provider override

If you want to use the local provider, you need to reference it by placing a file in your home folder under `~/.terraformrc` as follows:
```
provider_installation {

  dev_overrides {
      "registry.terraform.io/mezmo/mezmo" = "/Users/<YOUR USERNAME>/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Then, you can `plan` or `apply` a terraform files:

```bash
pushd examples/pipeline
terraform plan
popd
```

