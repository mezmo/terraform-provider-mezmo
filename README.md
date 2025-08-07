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

```hcl
provider_installation {
  dev_overrides {
      "registry.terraform.io/mezmo/mezmo" = "/Users/<YOUR USERNAME>/go/bin"
  }
}
```

Then, you can `plan` or `apply` a terraform files. This example assumes that there are `.tf`
files in `my-terraform-test` ready to be used.

Keep in mind if you have run the provider previously in your `my-terraform-test` directory before updating your config for dev_overrides your provider will be pinned to the remote repository version and will ignore your dev_overrides. To remedy this you will have to delete the .terraform in your `my-terraform-test` directory.

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

### ENV vars for Unit Tests

Optional environment variables can be provided on the test command line to display
additional debugging information when writing tests. Note that `fmt.Println` does not
work when running the provider directly, so these variables have no effect then. See
[Trace Logging During Execution](#trace-logging-during-execution-and-testing) for logging while using the provider.

- `DEBUG_ATTRIBUTES=1` - Displays the loaded state attributes when using the `StateHasExpectedValues` assertion

### Trace Logging During Execution and Testing

When running the provider directly or through integration tests, the APIs and their results can be printed to the screen.
For this, set `TF_LOG_PROVIDER_MEZMO=TRACE`, but be aware that it could print sensitive information.
This should only be used when debugging provider execution locally (including integration tests)!

#### Examples

- `-run` accepts a regex for the test name, and the path

- The path given shoulid match where the test file resides
- `TF_LOG_PROVIDER_MEZMO=TRACE` can be provided to see all api requests/responses/errors

#### Run all tests

```sh
make local-test
``````

#### Run a singular test

```sh
npm run local -- _TF_LOG_PROVIDER_MEZMO=TRACE go test -v -run 'TestAccAbsenceAlert_success' ./internal/provider/models/alerts/test
```

#### Run a group of tests

```sh
npm run local -- test -v -run 'TestAccChangeAlert.*_errors' go test -run ./internal/provider/models/alerts/test
```
