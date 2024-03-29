package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &RecordsetResource{}
var _ resource.ResourceWithImportState = &RecordsetResource{}

func NewRecordsetResource() resource.Resource {
	return &RecordsetResource{}
}

type RecordsetResource struct {
	client *powerdns.Client
}

type RecordsetResourceModel struct {
	Id       types.String `tfsdk:"id"`
	ZoneId   types.String `tfsdk:"zone_id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Ttl      types.Int64  `tfsdk:"ttl"`
	Records  types.List   `tfsdk:"records"`
}

func (r *RecordsetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_recordset"
}

func (t RecordsetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PowerDNS Zone",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "State ID for the record set (only needed for internal technical purposes).",
				Computed:            true,
			},
			"zone_id": schema.StringAttribute{
				MarkdownDescription: "ID of the zone this record set belongs to.",
				Required:            true,
			},
			"server_id": schema.StringAttribute{
				MarkdownDescription: "The id of the server.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name for record set (e.g. \"www.powerdns.com.\")",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of this record (e.g. \"A\", \"PTR\", \"MX\").",
				Required:            true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "DNS TTL of the records, in seconds.",
				Required:            true,
			},
			"records": schema.ListAttribute{
				MarkdownDescription: "All records in this record set.",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

func (r *RecordsetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*powerdns.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *powerdns.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r RecordsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordsetResourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	diags = RecordsetResourceModelToObject(ctx, data, recordset)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	serverId := data.ServerId.ValueString()
	recordSetName := data.Name.ValueString()
	recordSetType := data.Type.ValueString()
	recordSetTtl := data.Ttl.ValueInt64()
	tflog.Debug(ctx, "Creating record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       recordSetTtl,
		"records":   data.Records,
	})
	recordset, err := r.client.CreateRecordSet(ctx, serverId, zoneId, recordset)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to create record set '%s': %v",
				recordSetName,
				err))
		return
	}

	resp.Diagnostics.Append(recordsetObjectToResourceData(ctx, recordset, &data)...)
	tflog.Debug(ctx, "Created record set", map[string]interface{}{
		"id":        data.Id.ValueString(),
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       recordSetTtl,
		"records":   data.Records.Elements(),
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r RecordsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RecordsetResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	serverId := data.ServerId.ValueString()
	recordSetName := data.Name.ValueString()
	recordSetType := data.Type.ValueString()
	recordSetTtl := data.Ttl.ValueInt64()
	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
	})
	recordset, err := r.client.GetRecordSet(ctx, serverId, zoneId, recordSetName, recordSetType)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s' (type '%s'): %v", recordSetName, recordSetType, err))
		return
	}

	resp.Diagnostics.Append(recordsetObjectToResourceData(ctx, recordset, &data)...)
	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       recordSetTtl,
		"records":   data.Records.Elements(),
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r RecordsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RecordsetResourceModel

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	RecordsetResourceModelToObject(ctx, data, recordset)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	serverId := data.ServerId.ValueString()
	recordSetName := data.Name.ValueString()
	recordSetType := data.Type.ValueString()
	recordSetTtl := data.Ttl.ValueInt64()
	tflog.Debug(ctx, "Updating record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       recordSetTtl,
		"records":   data.Records.Elements(),
	})
	if err := r.client.UpdateRecordSet(ctx, serverId, zoneId, recordset); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update record set '%s': %v", recordset.Name, err))
		return
	}
	tflog.Debug(ctx, "Updated record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
	})

	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
	})
	recordset, err := r.client.GetRecordSet(ctx, serverId, zoneId, recordSetName, recordSetType)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s': %v", recordSetName, err))
		return
	}

	resp.Diagnostics.Append(recordsetObjectToResourceData(ctx, recordset, &data)...)
	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordSetName,
		"type":      recordSetType,
		"ttl":       recordSetTtl,
		"records":   data.Records.Elements(),
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r RecordsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RecordsetResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	RecordsetResourceModelToObject(ctx, data, recordset)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.ZoneId.ValueString()
	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Deleting record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordset.Name,
		"type":      recordset.Type,
	})
	if err := r.client.DeleteRecordSet(ctx, serverId, zoneId, recordset); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete record set '%s': %v", recordset.Name, err))
		return
	}
	tflog.Debug(ctx, "Deleted record set", map[string]interface{}{
		"zone_id":   zoneId,
		"server_id": serverId,
		"name":      recordset.Name,
		"type":      recordset.Type,
	})

	resp.State.RemoveResource(ctx)
}

func (r RecordsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splittedID := strings.Split(req.ID, "/")

	if len(splittedID) != 4 {
		resp.Diagnostics.AddError(
			"Resource Import ID invalid",
			fmt.Sprintf("ID '%s' should be in format 'server_id/zone_id/recordset_name/recordset_type'", req.ID),
		)
		return
	}
	serverID := splittedID[0]
	zoneID := splittedID[1]
	recordsetName := splittedID[2]
	recordsetType := splittedID[3]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_id"), serverID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), recordsetName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("type"), recordsetType)...)
}

func RecordsetResourceModelToObject(ctx context.Context, data RecordsetResourceModel, recordset *powerdns.RecordSet) diag.Diagnostics {
	var records []string
	diags := data.Records.ElementsAs(ctx, &records, false)
	if diags.HasError() {
		return diags
	}

	recordset.Name = data.Name.ValueString()
	recordset.Type = data.Type.ValueString()
	recordset.TTL = data.Ttl.ValueInt64()
	recordset.Records = records

	return diags
}

func recordsetObjectToResourceData(ctx context.Context, recordset *powerdns.RecordSet, data *RecordsetResourceModel) diag.Diagnostics {
	records := make([]attr.Value, len(recordset.Records))
	for i, record := range recordset.Records {
		records[i] = types.StringValue(record)
	}

	var diags diag.Diagnostics
	data.Id = types.StringValue(fmt.Sprintf("%s/%s/%s", data.ZoneId.ValueString(), data.Name.ValueString(), data.Type.ValueString()))
	data.Name = types.StringValue(recordset.Name)
	data.Type = types.StringValue(recordset.Type)
	data.Ttl = types.Int64Value(recordset.TTL)
	data.Records, diags = types.ListValue(types.StringType, records)

	return diags
}
