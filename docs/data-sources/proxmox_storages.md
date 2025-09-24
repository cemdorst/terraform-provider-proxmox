# proxmox_storages Data Source

The `proxmox_storages` data source allows you to retrieve information about all available storage configurations in your Proxmox VE cluster.

## Example Usage

```terraform
# Get all available storages
data "proxmox_storages" "all" {}

# Output all storages
output "all_storages" {
  value = data.proxmox_storages.all.storages
}

# Filter storages by type
output "dir_storages" {
  value = [for storage in data.proxmox_storages.all.storages : storage if storage.type == "dir"]
}

# Get storage names only
output "storage_names" {
  value = [for storage in data.proxmox_storages.all.storages : storage.storage]
}
```

Note: This data source requires proper Proxmox authentication. Configure your provider with:

```terraform
provider "proxmox" {
  endpoint     = "https://your-proxmox-host:8006"
  token_id     = "root@pam!terraform"
  token_secret = "your-secret-uuid"
  skip_verify  = true  # Only for development
}
```

## Schema

### Read-Only

- `id` (String) Data source identifier
- `storages` (List of Object) List of available storages

#### storages

- `storage` (String) Storage identifier/name
- `type` (String) Storage type (e.g., dir, lvm, nfs, cifs, etc.)
- `content` (String) Allowed content types (comma-separated list)
- `path` (String) Storage path (for local/directory storages)
- `priority` (Number) Storage priority
- `digest` (String) Storage configuration digest
- `prune_backups` (String) Prune backups configuration

## API Endpoint

This data source uses the Proxmox VE API endpoint: `/api2/json/storage/`
