package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &ZoneResource{}
var _ resource.ResourceWithImportState = &ZoneResource{}

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

type ZoneResource struct {
	client *powerdns.Client
}

type ZoneResourceModel struct {
	Id       types.String `tfsdk:"id"`
	ServerId types.String `tfsdk:"server_id"`
	Name     types.String `tfsdk:"name"`
	Kind     types.String `tfsdk:"kind"`
}

func (r *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (t *ZoneResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone := &powerdns.Zone{}
	zoneResourceDataToObject(ctx, data, zone)

	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Creating zone", map[string]interface{}{
		"server_id": serverId,
		"name":      zone.Name,
		"kind":      zone.Kind,
	})
	zone, err := r.client.CreateZone(ctx, serverId, zone)
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
		"id":        data.Id.ValueString(),
		"server_id": serverId,
		"name":      data.Name.ValueString(),
		"kind":      data.Kind.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var id string
	if !data.Id.IsUnknown() && !data.Id.IsNull() {
		id = data.Id.ValueString()
	} else {
		// If ID is not set, try to use name as id
		id = data.Name.ValueString()
	}

	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        id,
		"server_id": serverId,
	})
	zone, err := r.client.GetZone(ctx, serverId, id)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", id, err))
		return
	}

	zoneObjectToResourceData(ctx, zone, &data)
	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        data.Id.ValueString(),
		"server_id": serverId,
		"name":      data.Name.ValueString(),
		"kind":      data.Kind.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zone := &powerdns.Zone{}
	zoneResourceDataToObject(ctx, data, zone)

	var id string
	if !data.Id.IsUnknown() && !data.Id.IsNull() {
		id = data.Id.ValueString()
	} else {
		// If ID is not set, try to use name as id
		id = data.Name.ValueString()
	}

	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Updating zone", map[string]interface{}{
		"id":        id,
		"server_id": serverId,
		"name":      zone.Name,
		"kind":      zone.Kind,
		"dnssec":    zone.DNSSec,
		"masters":   zone.Masters,
	})
	if err := r.client.UpdateZone(ctx, serverId, id, zone); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update zone '%s': %v", id, err))
		return
	}
	tflog.Debug(ctx, "Updated zone", map[string]interface{}{
		"id":        id,
		"server_id": serverId,
	})

	tflog.Debug(ctx, "Reading zone", map[string]interface{}{
		"id":        id,
		"server_id": serverId,
	})
	zone, err := r.client.GetZone(ctx, serverId, id)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to get zone '%s': %v", id, err))
		return
	}

	zoneObjectToResourceData(ctx, zone, &data)
	tflog.Debug(ctx, "Read zone", map[string]interface{}{
		"id":        data.Id.ValueString(),
		"server_id": serverId,
		"name":      data.Name.ValueString(),
		"kind":      data.Kind.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := data.Id.ValueString()
	serverId := data.ServerId.ValueString()
	tflog.Debug(ctx, "Deleting zone", map[string]interface{}{
		"id":        zoneId,
		"server_id": serverId,
	})
	if err := r.client.DeleteZone(ctx, serverId, zoneId); err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete zone '%s': %v", zoneId, err))
		return
	}
	tflog.Debug(ctx, "Deleted zone", map[string]interface{}{
		"id":        zoneId,
		"server_id": serverId,
	})

	resp.State.RemoveResource(ctx)
}

func (r *ZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func zoneResourceDataToObject(ctx context.Context, data ZoneResourceModel, zone *powerdns.Zone) {
	zone.ID = data.Id.ValueString()
	zone.Name = data.Name.ValueString()
	zone.Kind = data.Kind.ValueString()
}

func zoneObjectToResourceData(ctx context.Context, zone *powerdns.Zone, data *ZoneResourceModel) {
	data.Id = types.StringValue(zone.ID)
	data.Name = types.StringValue(zone.Name)
	data.Kind = types.StringValue(zone.Kind)
}
