terraform {
  required_providers {
    proxmox = {
      source = "hashicorp/proxmox"
    }
  }
}

provider "proxmox" {
  endpoint     = "https://your-proxmox-host:8006"
  token_id     = "root@pam!terraform"
  token_secret = "your-secret-uuid"
  skip_verify  = true
}

# Get all available storages
data "proxmox_storages" "all" {}

# Output all storages
output "all_storages" {
  value = data.proxmox_storages.all.storages
}

# Output only storage names
output "storage_names" {
  value = [for storage in data.proxmox_storages.all.storages : storage.storage]
}

# Filter storages by type
output "dir_storages" {
  value = [for storage in data.proxmox_storages.all.storages : storage if storage.type == "dir"]
}
