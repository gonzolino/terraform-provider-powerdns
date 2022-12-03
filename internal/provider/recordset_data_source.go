package provider

import (
	"context"
	"fmt"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &RecordsetDataSource{}

func NewRecordsetDataSource() datasource.DataSource {
	return &RecordsetDataSource{}
}

// RecordsetDataSource defines the data source implementation.
type RecordsetDataSource struct {
	client *powerdns.Client
}

// RecordsetDataSourceModel describes the data source data model.
type RecordsetDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	ZoneId   types.String `tfsdk:"zone_id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Ttl      types.Int64  `tfsdk:"ttl"`
	Records  types.List   `tfsdk:"records"`
}

func (d *RecordsetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recordset"
}

func (d RecordsetDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "PowerDNS Resource Record Set",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "State ID for the record set (only needed for internal technical purposes).",
				Type:                types.StringType,
				Computed:            true,
			},
			"zone_id": {
				MarkdownDescription: "ID of the zone this record set belongs to.",
				Type:                types.StringType,
				Required:            true,
			},
			"server_id": {
				MarkdownDescription: "The id of the server.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "Name for record set (e.g. \"www.powerdns.com.\")",
				Type:                types.StringType,
				Required:            true,
			},
			"type": {
				MarkdownDescription: "Type of this record (e.g. \"A\", \"PTR\", \"MX\"). Required if record name is not unique.",
				Type:                types.StringType,
				Required:            true,
			},
			"ttl": {
				MarkdownDescription: " DNS TTL of the records, in seconds.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"records": {
				MarkdownDescription: "All records in this record set.",
				Type:                types.ListType{ElemType: types.StringType},
				Computed:            true,
			},
		},
	}, nil
}

func (d *RecordsetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d RecordsetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordsetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	serverId := data.ServerId.ValueString()
	recordSetName := data.Name.ValueString()
	recordSetType := data.Type.ValueString()
	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
	})
	recordset, err := d.client.GetRecordSet(ctx, serverId, zoneId, recordSetName, recordSetType)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s' (type '%s'): %v", recordSetName, recordSetType, err))
		return
	}

	records := make([]attr.Value, len(recordset.Records))
	for i, record := range recordset.Records {
		records[i] = types.StringValue(record)
	}

	var diags diag.Diagnostics
	data.Id = types.StringValue(fmt.Sprintf("%s/%s/%s", zoneId, recordSetName, recordSetType))
	data.Name = types.StringValue(recordset.Name)
	data.Type = types.StringValue(recordset.Type)
	data.Ttl = types.Int64Value(recordset.TTL)
	data.Records, diags = types.ListValue(types.StringType, records)

	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       data.Ttl.ValueInt64(),
		"records":   data.Records,
	})

	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
