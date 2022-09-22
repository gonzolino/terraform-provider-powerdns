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

	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
	})
	recordset, err := d.client.GetRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, data.Name.Value, data.Type.Value)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s' (type '%s'): %v", data.Name.Value, data.Type.Value, err))
		return
	}

	records := make([]attr.Value, len(recordset.Records))
	for i, record := range recordset.Records {
		records[i] = types.String{Value: record}
	}

	data.Id = types.String{Value: fmt.Sprintf("%s/%s/%s", data.ZoneId.Value, data.Name.Value, data.Type.Value)}
	data.Name = types.String{Value: recordset.Name}
	data.Type = types.String{Value: recordset.Type}
	data.Ttl = types.Int64{Value: recordset.TTL}
	data.Records = types.List{
		ElemType: types.StringType,
		Elems:    records,
	}

	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
