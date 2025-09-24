// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoragesDataSource{}

func NewStoragesDataSource() datasource.DataSource {
	return &StoragesDataSource{}
}

// StoragesDataSource defines the data source implementation.
type StoragesDataSource struct {
	client *ProxmoxClient
}

// StoragesDataSourceModel describes the data source data model.
type StoragesDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Storages []StorageModel `tfsdk:"storages"`
}

// StorageModel describes a single storage entry.
type StorageModel struct {
	Storage      types.String `tfsdk:"storage"`
	Type         types.String `tfsdk:"type"`
	Content      types.String `tfsdk:"content"`
	Path         types.String `tfsdk:"path"`
	Priority     types.Int64  `tfsdk:"priority"`
	Digest       types.String `tfsdk:"digest"`
	PruneBackups types.String `tfsdk:"prune_backups"`
}

func (d *StoragesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storages"
}

func (d *StoragesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all available Proxmox VE storages.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier",
				Computed:            true,
			},
			"storages": schema.ListNestedAttribute{
				MarkdownDescription: "List of available storages",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"storage": schema.StringAttribute{
							MarkdownDescription: "Storage identifier",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Storage type (e.g., dir, lvm, nfs, etc.)",
							Computed:            true,
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "Allowed content types",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Storage path",
							Computed:            true,
						},
						"priority": schema.Int64Attribute{
							MarkdownDescription: "Storage priority",
							Computed:            true,
						},
						"digest": schema.StringAttribute{
							MarkdownDescription: "Storage digest",
							Computed:            true,
						},
						"prune_backups": schema.StringAttribute{
							MarkdownDescription: "Prune backups configuration",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *StoragesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ProxmoxClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProxmoxClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *StoragesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StoragesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Proxmox storages")

	// Make API request to get storages
	httpResp, err := d.client.DoRequest("GET", "/storage", nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read storages, got error: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to read storages, got status %d: %s", httpResp.StatusCode, string(body)),
		)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read response body: %s", err))
		return
	}

	var storageResponse struct {
		Data []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &storageResponse); err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
		return
	}

	// Convert response to model
	storages := make([]StorageModel, len(storageResponse.Data))
	for i, storageData := range storageResponse.Data {
		storage := StorageModel{}

		if val, ok := storageData["storage"].(string); ok {
			storage.Storage = types.StringValue(val)
		} else {
			storage.Storage = types.StringNull()
		}

		if val, ok := storageData["type"].(string); ok {
			storage.Type = types.StringValue(val)
		} else {
			storage.Type = types.StringNull()
		}

		if val, ok := storageData["content"].(string); ok {
			storage.Content = types.StringValue(val)
		} else {
			storage.Content = types.StringNull()
		}

		if val, ok := storageData["path"].(string); ok {
			storage.Path = types.StringValue(val)
		} else {
			storage.Path = types.StringNull()
		}

		if val, ok := storageData["priority"].(float64); ok {
			storage.Priority = types.Int64Value(int64(val))
		} else {
			storage.Priority = types.Int64Null()
		}

		if val, ok := storageData["digest"].(string); ok {
			storage.Digest = types.StringValue(val)
		} else {
			storage.Digest = types.StringNull()
		}

		if val, ok := storageData["prune-backups"].(string); ok {
			storage.PruneBackups = types.StringValue(val)
		} else {
			storage.PruneBackups = types.StringNull()
		}

		storages[i] = storage
	}

	data.Storages = storages
	data.ID = types.StringValue("storages")

	tflog.Debug(ctx, fmt.Sprintf("Found %d storages", len(storages)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
