package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type recordsetResourceType struct{}

func (t recordsetResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (t recordsetResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return recordsetResource{
		provider: provider,
	}, diags
}

type recordsetResourceData struct {
	Id       types.String `tfsdk:"id"`
	ZoneId   types.String `tfsdk:"zone_id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Ttl      types.Int64  `tfsdk:"ttl"`
	Records  types.List   `tfsdk:"records"`
}

type recordsetResource struct {
	provider provider
}

func (r recordsetResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data recordsetResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	diags = recordsetResourceDataToObject(ctx, data, recordset)
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
	recordset, err := r.provider.client.CreateRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset)
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

func (r recordsetResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data recordsetResourceData

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
	recordset, err := r.provider.client.GetRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, data.Name.Value, data.Type.Value)
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

func (r recordsetResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data recordsetResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	recordsetResourceDataToObject(ctx, data, recordset)
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
	if err := r.provider.client.UpdateRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset); err != nil {
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
	recordset, err := r.provider.client.GetRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, data.Name.Value, data.Type.Value)
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

func (r recordsetResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data recordsetResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordset := &powerdns.RecordSet{}
	recordsetResourceDataToObject(ctx, data, recordset)
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
	if err := r.provider.client.DeleteRecordSet(ctx, data.ServerId.Value, data.ZoneId.Value, recordset); err != nil {
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

func (r recordsetResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("server_id"), serverID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("zone_id"), zoneID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("name"), recordsetName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("type"), recordsetType)...)
}

func recordsetResourceDataToObject(ctx context.Context, data recordsetResourceData, recordset *powerdns.RecordSet) diag.Diagnostics {
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

func recordsetObjectToResourceData(ctx context.Context, recordset *powerdns.RecordSet, data *recordsetResourceData) {
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
