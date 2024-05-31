# Mezmo Terraform Provider

The Mezmo Terraform Provider allows organizations to manage Pipelines (sources, processors and destinations)
programmatically via Terraform.

You can download this repo to create your own local provider or you can use the Hashicorp registry.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building the Provider Locally
Building the provider will build the provider binary in the project root. It's not very useful in most
cases, so see the instructions in [Running the Provider Locally](#running-the-provider-locally).
```shell
go build ./...
```

## Running the Provider Locally

To build and install the provider locally, run `go install .`. This will build the provider and put the provider
binary in the `$GOPATH/bin` directory. You will reference this location in the next step.

### Adding the Provider Override

If you want to use the local provider, you need to reference it by placing a file in your `$HOME` directory called `.terraformrc`.
This setting tells terraform to override the remote registry where the provider is usually downloaded from in favor of a local directory.
The value of this setting should be your `$GOPATH`, which is where the installed binary gets placed.
This value is usually `$HOME/go/bin`.
```
provider_installation {
  dev_overrides {
      "registry.terraform.io/mezmo/mezmo" = "/Users/<YOUR USERNAME>/go/bin"
  }
}
```

Then, you can `plan` or `apply` a terraform files. This example assumes that there are `.tf`
files in `my-terraform-test` ready to be used.

```bash
cd my-terraform-test
terraform init
terraform plan
terraform apply
```

## Generating the Docs

When schemas are changed (descriptions, types) during development, the documentation for the components must be re-generated.
To do this, run `go generate` to make sure all changes are documented.

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

