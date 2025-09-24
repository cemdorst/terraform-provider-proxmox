// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStoragesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccStoragesDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.proxmox_storages.test", "id", "storages"),
					resource.TestCheckResourceAttrSet("data.proxmox_storages.test", "storages.#"),
				),
			},
		},
	})
}

func testAccStoragesDataSourceConfig() string {
	return fmt.Sprintf(`
provider "proxmox" {
  endpoint     = "%s"
  token_id     = "%s"
  token_secret = "%s"
  skip_verify  = true
}

data "proxmox_storages" "test" {}
`, testEndpoint(), testTokenID(), testTokenSecret())
}

func testEndpoint() string {
	endpoint := os.Getenv("PROXMOX_ENDPOINT")
	if endpoint == "" {
		return "https://proxmox.example.com:8006"
	}
	return endpoint
}

func testTokenID() string {
	tokenID := os.Getenv("PROXMOX_TOKEN_ID")
	if tokenID == "" {
		return "root@pam!test"
	}
	return tokenID
}

func testTokenSecret() string {
	tokenSecret := os.Getenv("PROXMOX_TOKEN_SECRET")
	if tokenSecret == "" {
		return "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
	return tokenSecret
}
