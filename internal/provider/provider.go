// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ProxmoxProvider satisfies various provider interfaces.
var _ provider.Provider = &ProxmoxProvider{}

// ProxmoxClient wraps the HTTP client for Proxmox API communication.
type ProxmoxClient struct {
	HTTPClient  *http.Client
	Endpoint    string
	TokenID     string
	TokenSecret string
}

// DoRequest makes an HTTP request to the Proxmox API.
func (c *ProxmoxClient) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	url := strings.TrimSuffix(c.Endpoint, "/") + "/api2/json" + path
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "PVEAPIToken="+c.TokenID+"="+c.TokenSecret)
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// ProxmoxProvider defines the provider implementation.
type ProxmoxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProxmoxProviderModel describes the provider data model.
type ProxmoxProviderModel struct {
	Endpoint    types.String `tfsdk:"endpoint"`
	TokenID     types.String `tfsdk:"token_id"`
	TokenSecret types.String `tfsdk:"token_secret"`
	SkipVerify  types.Bool   `tfsdk:"skip_verify"`
}

func (p *ProxmoxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "proxmox"
	resp.Version = p.version
}

func (p *ProxmoxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Proxmox API endpoint URL (e.g., https://proxmox.example.com:8006)",
				Required:            true,
			},
			"token_id": schema.StringAttribute{
				MarkdownDescription: "Proxmox API token ID (e.g., root@pam!mytesttoken)",
				Required:            true,
			},
			"token_secret": schema.StringAttribute{
				MarkdownDescription: "Proxmox API token secret",
				Required:            true,
				Sensitive:           true,
			},
			"skip_verify": schema.BoolAttribute{
				MarkdownDescription: "Skip TLS certificate verification",
				Optional:            true,
			},
		},
	}
}

func (p *ProxmoxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProxmoxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"The provider cannot create the Proxmox API client as there is a missing or empty value for the Proxmox endpoint.",
		)
		return
	}

	if data.TokenID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"The provider cannot create the Proxmox API client as there is a missing or empty value for the Proxmox API token ID.",
		)
		return
	}

	if data.TokenSecret.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"The provider cannot create the Proxmox API client as there is a missing or empty value for the Proxmox API token secret.",
		)
		return
	}

	// Validate token ID format
	tokenID := data.TokenID.ValueString()
	if !strings.Contains(tokenID, "!") {
		resp.Diagnostics.AddError(
			"Invalid Token ID Format",
			"The API token ID should contain a '!' character and follow the format 'user@realm!tokenname' (e.g., 'root@pam!mytesttoken').",
		)
		return
	}

	// Create HTTP client with optional TLS skip verification
	transport := &http.Transport{}
	if !data.SkipVerify.IsNull() && data.SkipVerify.ValueBool() {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &ProxmoxClient{
		HTTPClient:  &http.Client{Transport: transport},
		Endpoint:    data.Endpoint.ValueString(),
		TokenID:     data.TokenID.ValueString(),
		TokenSecret: data.TokenSecret.ValueString(),
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ProxmoxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ProxmoxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewStoragesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ProxmoxProvider{
			version: version,
		}
	}
}
