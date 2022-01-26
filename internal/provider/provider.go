package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// provider satisfies the tfsdk.Provider interface and usually is included
// with all Resource and DataSource implementations.
type provider struct {
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

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
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

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"powerdns_zone": zoneResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"powerdns_zone": zoneDataSourceType{},
	}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "PowerDNS API key for authentication",
				Optional:            true,
				Sensitive:           true,
				Type:                types.StringType,
			},
			"server_url": {
				MarkdownDescription: "PowerDNS server URL",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
