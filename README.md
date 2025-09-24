# Terraform Provider for Proxmox

This Terraform provider allows you to manage Proxmox VE resources using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Features

- **Storage Discovery**: List and query all available storage configurations in your Proxmox cluster
- **Token-based Authentication**: Secure API token authentication with Proxmox VE
- **TLS Configuration**: Optional TLS certificate verification skipping for development environments

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- Proxmox VE 6.x or later
- Proxmox API token with appropriate permissions

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Using the Provider

### Authentication

You'll need to create an API token in Proxmox VE:

1. Log into your Proxmox web interface
2. Go to Datacenter -> Permissions -> API Tokens
3. Create a new token with appropriate privileges

### Basic Configuration

```hcl
terraform {
  required_providers {
    proxmox = {
      source = "cemdorst/proxmox"
    }
  }
}

provider "proxmox" {
  endpoint     = "https://your-proxmox-server.example.com:8006"
  token_id     = "root@pam!terraform"
  token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  skip_verify  = true  # Only for development/testing
}
```

### Environment Variables

You can also use environment variables:

- `PROXMOX_ENDPOINT` - Proxmox API endpoint URL
- `PROXMOX_TOKEN_ID` - API token ID (e.g., root@pam!terraform)
- `PROXMOX_TOKEN_SECRET` - API token secret UUID

### Storages Data Source

```hcl
# Get all available storages
data "proxmox_storages" "all" {}

# Output storage names
output "storage_names" {
  value = [for storage in data.proxmox_storages.all.storages : storage.storage]
}
```

## Provider Configuration

### Arguments

- `endpoint` (String, Required) - The Proxmox API endpoint URL (e.g., `https://proxmox.example.com:8006`)
- `token_id` (String, Required) - The Proxmox API token ID (e.g., `root@pam!terraform`)
- `token_secret` (String, Required, Sensitive) - The Proxmox API token secret UUID
- `skip_verify` (Boolean, Optional) - Skip TLS certificate verification (default: false)

## Data Sources

### `proxmox_storages`

Retrieves information about all available storage configurations in your Proxmox VE cluster.

#### Example Usage

```hcl
# Get all available storages
data "proxmox_storages" "all" {}

# Output storage names
output "storage_names" {
  value = [for storage in data.proxmox_storages.all.storages : storage.storage]
}

# Filter storages by type
output "dir_storages" {
  value = [for storage in data.proxmox_storages.all.storages : storage if storage.type == "dir"]
}
```

#### Attributes Reference

- `id` (String) - Data source identifier
- `storages` (List of Object) - List of available storages
  - `storage` (String) - Storage identifier/name
  - `type` (String) - Storage type (e.g., dir, lvm, nfs, cifs, etc.)
  - `content` (String) - Allowed content types (comma-separated list)
  - `path` (String) - Storage path (for local/directory storages)
  - `priority` (Number) - Storage priority
  - `digest` (String) - Storage configuration digest
  - `prune_backups` (String) - Prune backups configuration

## Developing the Provider

### Prerequisites

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

### Building

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

### Testing

To run the unit tests:

```shell
go test ./internal/provider/...
```

To run the acceptance tests (requires a real Proxmox environment):

```shell
export PROXMOX_ENDPOINT="https://your-proxmox-server.example.com:8006"
export PROXMOX_TOKEN_ID="root@pam!terraform"
export PROXMOX_TOKEN_SECRET="your-secret-uuid"
export TF_ACC=1
go test ./internal/provider/... -v
```

*Note:* Acceptance tests create real resources on your Proxmox server.

### Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

To add a new dependency:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Mozilla Public License v2.0 - see the LICENSE file for details.