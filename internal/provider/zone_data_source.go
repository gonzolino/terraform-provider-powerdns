package provider

import (
	"context"
	"fmt"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ZoneDataSource{}

func NewZoneDataSource() datasource.DataSource {
	return &ZoneDataSource{}
}

// ZoneDataSource defines the data source implementation.
type ZoneDataSource struct {
	client *powerdns.Client
}

// ZoneDataSourceModel describes the data source data model.
type ZoneDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Kind     types.String `tfsdk:"kind"`
}

func (d *ZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (d ZoneDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (d *ZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*powerdns.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *powerdns.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d ZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ZoneDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
	})
	zone, err := d.client.GetZone(ctx, data.ServerId.Value, data.Id.Value)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", data.Id.Value, err))
		return
	}

	data.Id = types.String{Value: zone.ID}
	data.Name = types.String{Value: zone.Name}
	data.Kind = types.String{Value: zone.Kind}

	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        zone.ID,
		"server_id": data.ServerId.Value,
		"name":      zone.Name,
		"kind":      zone.Kind,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
