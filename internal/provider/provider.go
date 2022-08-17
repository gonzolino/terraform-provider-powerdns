package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type powerdnsProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.
	client *powerdns.Client

	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	APIKey    types.String `tfsdk:"api_key"`
	ServerURL types.String `tfsdk:"server_url"`
}

func (p *powerdnsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var apiKey string
	if data.APIKey.Unknown {
		resp.Diagnostics.AddWarning("API Key is not set", "API Key is not set. This is required for authentication.")
		return
	}
	if data.APIKey.Null {
		apiKey = os.Getenv("POWERDNS_API_KEY")
	} else {
		apiKey = data.APIKey.Value
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("API Key is not set", "API Key is not set. This is required for authentication.")
		return
	}

	var serverURL string
	if data.ServerURL.Unknown {
		resp.Diagnostics.AddWarning("Server URL is not set", "Server URL is not set. Can't connect to PowerDNS API.")
		return
	}
	if data.ServerURL.Null {
		serverURL = os.Getenv("POWERDNS_SERVER_URL")
	} else {
		serverURL = data.ServerURL.Value
	}

	// Configuration values are now available.
	parsedServerURL, err := url.Parse(serverURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Server URL",
			fmt.Sprintf("Invalid server URL: %v", err),
		)
		return
	}

	p.client = powerdns.New(ctx, apiKey, parsedServerURL.Host, parsedServerURL.Path, parsedServerURL.Scheme)
	p.configured = true
}

func (p *powerdnsProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"powerdns_recordset": recordsetResourceType{},
		"powerdns_zone":      zoneResourceType{},
	}, nil
}

func (p *powerdnsProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"powerdns_recordset": recordsetDataSourceType{},
		"powerdns_zone":      zoneDataSourceType{},
	}, nil
}

func (p *powerdnsProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "The PowerDNS provider allows modifying zone content and metadata using the PowerDNS API.",
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "PowerDNS API key for authentication. Can be set via environment variable `POWERDNS_API_KEY`.",
				Optional:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
			"server_url": {
				MarkdownDescription: "PowerDNS server URL. Can be set via environment variable `POWERDNS_SERVER_URL`.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &powerdnsProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (powerdnsProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*powerdnsProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return powerdnsProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return powerdnsProvider{}, diags
	}

	return *p, diags
}
