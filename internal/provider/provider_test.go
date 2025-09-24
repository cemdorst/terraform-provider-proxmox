// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"proxmox": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PROXMOX_ENDPOINT") == "" {
		t.Skip("PROXMOX_ENDPOINT environment variable must be set for acceptance tests")
	}
	if os.Getenv("PROXMOX_TOKEN_ID") == "" {
		t.Skip("PROXMOX_TOKEN_ID environment variable must be set for acceptance tests")
	}
	if os.Getenv("PROXMOX_TOKEN_SECRET") == "" {
		t.Skip("PROXMOX_TOKEN_SECRET environment variable must be set for acceptance tests")
	}
}
