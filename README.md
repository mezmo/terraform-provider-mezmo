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

## Testing

To run the full test suite, which will start Docker containers for the required services, run:

```sh
make test
```

### Testing Locally

Docker containers for the required services can be started, then tests can run
individually against them. The services will run in the foreground, which will display
their logs.  Wait for `pipeline-service` to be up and running, indicated when the log
shows the loaded API routes. In a separate window, run:

```sh
make start
```
Each test command will be preceeded by loading the proper environment variables. Depending
on the noted shell, that is done by:

```sh
shell|bash> env $(cat env/local.env) #...rest of command
fish> env (cat env/local.env) #...rest of command
```

### ENV vars
Optional environment variables can be provided on the test command line to display
additional debugging information.

* `DEBUG_SOURCE=1` - Displays the API request/responses for sources
* `DEBUG_PROCESSOR=1` - Displays the API request/responses for processors
* `DEBUG_DESTINATION=1` - Displays the API request/responses for destinations
* `DEBUG_ALERT=1` - Displays the API request/responses for alerts

* `DEBUG_ATTRIBUTES=1` - Displays the loaded state attributes when using the `StateHasExpectedValues` assertion

#### Examples
* `-run` accepts a regex for the test name, and the path
* The path given shoulid match where the test file resides

**Run all tests**

```sh
make local-test
``````
**Run a singular test**

```sh
env $(cat env/local.env) DEBUG_ATTRIBUTES=1 DEBUG_ALERT=1 go test -v -run 'TestAbsenceAlert_success' ./internal/provider/models/alerts/test
```

**Run a group of tests**

```sh
env $(cat env/local.env) test -v -run 'TestChangeAlert.*_errors' ./internal/provider/models/alerts/test
```

