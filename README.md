# Terraform Provider for [incident.io](https://incident.io)

This Terraform provider helps you to configure your [incident.io](https://incident.io) account.

1. Get an API key in https://app.incident.io/settings/api-keys
2. Configure https://incident.io using the provider:

```hcl
terraform {
  required_providers {
    incidentio = {
      source = "multani/incidentio"
    }
  }
}

variable "incidentio_api_key" {
  description = <<EOF
An incident.io API key.

Get one at https://app.incident.io/settings/api-keys
EOF
}

provider "incidentio" {
  api_key = var.incidentio_api_key
}

resource "incidentio_incident_role" "spectator" {
  name        = "Spectator"
  short_form  = "spectator"
  description = "A person that enjoys eating popcorn when things are burning."

  instructions = <<EOF
- Grab some popcorn
- Silently watch the incident happening
- Congrat the other roles once the incident has been resolved
EOF

  # This role is mostly probably not required, most of the time.
  required   = false
}
```

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
