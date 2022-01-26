package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type zoneDataSourceType struct{}

func (t zoneDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "PowerDNS Zone",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Opaque zone id, assigned by the server.",
				Type:                types.StringType,
				Required:            true,
			},
			"server_id": {
				MarkdownDescription: "The id of the server.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "Name of the zone (e.g. \"example.com.\") MUST have a trailing dot.",
				Type:                types.StringType,
				Computed:            true,
			},
			"kind": {
				MarkdownDescription: "Zone kind, one of \"Native\", \"Master\", \"Slave\".",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (t zoneDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return zoneDataSource{
		provider: provider,
	}, diags
}

type zoneDataSourceData struct {
	Id       types.String `tfsdk:"id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Kind     types.String `tfsdk:"kind"`
}

type zoneDataSource struct {
	provider provider
}

func (d zoneDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data zoneDataSourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading zone", "id", data.Id.Value, "server_id", data.ServerId.Value)
	zone, err := d.provider.client.GetZone(ctx, data.ServerId.Value, data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", data.Id.Value, err))
		return
	}

	data.Id = types.String{Value: zone.ID}
	data.Name = types.String{Value: zone.Name}
	data.Kind = types.String{Value: zone.Kind}

	tflog.Debug(ctx, "Read zone", "id", zone.ID, "server_id", data.ServerId.Value, "name", zone.Name, "kind", zone.Kind)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
