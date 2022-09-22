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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func (t RecordsetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "PowerDNS Zone",

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
				MarkdownDescription: "Type of this record (e.g. \"A\", \"PTR\", \"MX\").",
				Type:                types.StringType,
				Required:            true,
			},
			"ttl": {
				MarkdownDescription: "DNS TTL of the records, in seconds.",
				Type:                types.Int64Type,
				Required:            true,
			},
			"records": {
				MarkdownDescription: "All records in this record set.",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            true,
			},
		},
	}, nil
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

	tflog.Debug(ctx, "Creating record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records,
	})
	recordset, err := r.client.CreateRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to create record set '%s': %v",
				data.Name.Value,
				err))
		return
	}

	recordsetObjectToResourceData(ctx, recordset, &data)
	tflog.Debug(ctx, "Created record set", map[string]interface{}{
		"id":        data.Id.Value,
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records.Elems,
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

	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
	})
	recordset, err := r.client.GetRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, data.Name.Value, data.Type.Value)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s' (type '%s'): %v", data.Name.Value, data.Type.Value, err))
		return
	}

	recordsetObjectToResourceData(ctx, recordset, &data)
	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records.Elems,
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

	tflog.Debug(ctx, "Updating record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records.Elems,
	})
	if err := r.client.UpdateRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update record set '%s': %v", recordset.Name, err))
		return
	}
	tflog.Debug(ctx, "Updated record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
	})

	tflog.Debug(ctx, "Reading record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
	})
	recordset, err := r.client.GetRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, data.Name.Value, data.Type.Value)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get record set '%s': %v", data.Name.Value, err))
		return
	}

	recordsetObjectToResourceData(ctx, recordset, &data)
	tflog.Debug(ctx, "Read record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"type":      data.Type.Value,
		"ttl":       data.Ttl.Value,
		"records":   data.Records.Elems,
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

	tflog.Debug(ctx, "Deleting record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
		"name":      recordset.Name,
		"type":      recordset.Type,
	})
	if err := r.client.DeleteRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete record set '%s': %v", recordset.Name, err))
		return
	}
	tflog.Debug(ctx, "Deleted record set", map[string]interface{}{
		"zone_id":   data.ZoneId.Value,
		"server_id": data.ServerId.Value,
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

	recordset.Name = data.Name.Value
	recordset.Type = data.Type.Value
	recordset.TTL = data.Ttl.Value
	recordset.Records = records

	return diags
}

func recordsetObjectToResourceData(ctx context.Context, recordset *powerdns.RecordSet, data *RecordsetResourceModel) {
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
}
