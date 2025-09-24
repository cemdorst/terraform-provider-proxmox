terraform {
  required_providers {
    proxmox = {
      source = "cemdorst/proxmox"
    }
  }
}

provider "proxmox" {
  endpoint     = "https://proxmox.example.com:8006"
  token_id     = "root@pam!terraform"
  token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  skip_verify  = true
}
