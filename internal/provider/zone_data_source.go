package provider

import (
	"context"
	"fmt"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

func (d ZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PowerDNS Zone",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Opaque zone id, assigned by the server.",
				Required:            true,
			},
			"server_id": schema.StringAttribute{
				MarkdownDescription: "The id of the server.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the zone (e.g. \"example.com.\") MUST have a trailing dot.",
				Computed:            true,
			},
			"kind": schema.StringAttribute{
				MarkdownDescription: "Zone kind, one of \"Native\", \"Master\", \"Slave\".",
				Computed:            true,
			},
		},
	}
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

	zoneId := data.Id.ValueString()
	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        zoneId,
		"server_id": serverId,
	})
	zone, err := d.client.GetZone(ctx, serverId, zoneId)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", zoneId, err))
		return
	}

	data.Id = types.StringValue(zone.ID)
	data.Name = types.StringValue(zone.Name)
	data.Kind = types.StringValue(zone.Kind)

	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        zone.ID,
		"server_id": serverId,
		"name":      zone.Name,
		"kind":      zone.Kind,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
