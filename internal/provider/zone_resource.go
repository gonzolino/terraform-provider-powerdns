package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type zoneResourceType struct{}

func (t zoneResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "PowerDNS Zone",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Opaque zone id, assigned by the server.",
				Type:                types.StringType,
				Computed:            true,
			},
			"server_id": {
				MarkdownDescription: "The id of the server.",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "Name of the zone (e.g. \"example.com.\") MUST have a trailing dot.",
				Type:                types.StringType,
				Required:            true,
			},
			"kind": {
				MarkdownDescription: "Zone kind, one of \"Native\", \"Master\", \"Slave\".",
				Type:                types.StringType,
				Required:            true,
			},
		},
	}, nil
}

func (t zoneResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return zoneResource{
		provider: provider,
	}, diags
}

type zoneResourceData struct {
	Id       types.String `tfsdk:"id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Kind     types.String `tfsdk:"kind"`
}

type zoneResource struct {
	provider provider
}

func (r zoneResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data zoneResourceData

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone := &powerdns.Zone{}
	zoneResourceDataToObject(ctx, data, zone)

	tflog.Debug(ctx, "Creating zone", map[string]interface{}{
		"server_id": data.ServerId.Value,
		"name":      zone.Name,
		"kind":      zone.Kind,
	})
	zone, err := r.provider.client.CreateZone(ctx, data.ServerId.Value, zone)
	if err != nil {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to create zone '%s': %v",
				zone.Name,
				err))
		return
	}

	zoneObjectToResourceData(ctx, zone, &data)
	tflog.Debug(ctx, "Created zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"kind":      data.Kind.Value,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r zoneResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data zoneResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var id string
	if !data.Id.Unknown && !data.Id.Null {
		id = data.Id.Value
	} else {
		// If ID is not set, try to use name as id
		id = data.Name.Value
	}

	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        id,
		"server_id": data.ServerId.Value,
	})
	zone, err := r.provider.client.GetZone(ctx, data.ServerId.Value, id)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", id, err))
		return
	}

	zoneObjectToResourceData(ctx, zone, &data)
	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"kind":      data.Kind.Value,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r zoneResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data zoneResourceData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone := &powerdns.Zone{}
	zoneResourceDataToObject(ctx, data, zone)

	var id string
	if !data.Id.Unknown && !data.Id.Null {
		id = data.Id.Value
	} else {
		// If ID is not set, try to use name as id
		id = data.Name.Value
	}

	tflog.Debug(ctx, "Updating zone", map[string]interface{}{
		"id":        id,
		"server_id": data.ServerId.Value,
		"name":      zone.Name,
		"kind":      zone.Kind,
		"dnssec":    zone.DNSSec,
		"masters":   zone.Masters,
	})
	if err := r.provider.client.UpdateZone(ctx, data.ServerId.Value, id, zone); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update zone '%s': %v", id, err))
		return
	}
	tflog.Debug(ctx, "Updated zone", map[string]interface{}{
		"id":        id,
		"server_id": data.ServerId.Value,
	})

	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        id,
		"server_id": data.ServerId.Value,
	})
	zone, err := r.provider.client.GetZone(ctx, data.ServerId.Value, id)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", id, err))
		return
	}

	zoneObjectToResourceData(ctx, zone, &data)
	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
		"name":      data.Name.Value,
		"kind":      data.Kind.Value,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r zoneResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data zoneResourceData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
	})
	if err := r.provider.client.DeleteZone(ctx, data.ServerId.Value, data.Id.Value); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete zone '%s': %v", data.Id.Value, err))
		return
	}
	tflog.Debug(ctx, "Deleted zone", map[string]interface{}{
		"id":        data.Id.Value,
		"server_id": data.ServerId.Value,
	})

	resp.State.RemoveResource(ctx)
}

func (r zoneResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	splittedID := strings.Split(req.ID, "/")

	if len(splittedID) != 2 {
		resp.Diagnostics.AddError(
			"Resource Import ID invalid",
			fmt.Sprintf("ID '%s' should be in format 'server_id/zone_id'", req.ID),
		)
		return
	}
	serverID := splittedID[0]
	zoneID := splittedID[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_id"), serverID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), zoneID)...)

	// tfsdk.ResourceImportStateNotImplemented(ctx, "", resp)
	// tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func zoneResourceDataToObject(ctx context.Context, data zoneResourceData, zone *powerdns.Zone) {
	zone.ID = data.Id.Value
	zone.Name = data.Name.Value
	zone.Kind = data.Kind.Value
}

func zoneObjectToResourceData(ctx context.Context, zone *powerdns.Zone, data *zoneResourceData) {
	data.Id = types.String{Value: zone.ID}
	data.Name = types.String{Value: zone.Name}
	data.Kind = types.String{Value: zone.Kind}
}
