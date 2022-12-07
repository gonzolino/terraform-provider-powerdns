package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/gonzolino/terraform-provider-powerdns/internal/powerdns"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure PowerdnsProvider satisfies various provider interfaces.
var _ provider.Provider = &PowerdnsProvider{}
var _ provider.ProviderWithMetadata = &PowerdnsProvider{}

// PowerdnsProvider defines the provider implementation.
type PowerdnsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PowerdnsProviderModel describes the provider data model.
type PowerdnsProviderModel struct {
	APIKey    types.String `tfsdk:"api_key"`
	ServerURL types.String `tfsdk:"server_url"`
}

func (p *PowerdnsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "powerdns"
	resp.Version = p.version
}

func (p *PowerdnsProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The PowerDNS provider allows modifying zone content and metadata using the PowerDNS API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "PowerDNS API key for authentication. Can be set via environment variable `POWERDNS_API_KEY`.",
				Optional:            true,
				Sensitive:           true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "PowerDNS server URL. Can be set via environment variable `POWERDNS_SERVER_URL`.",
				Optional:            true,
			},
		},
	}
}

func (p *PowerdnsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PowerdnsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var apiKey string
	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddWarning("API Key is not set", "API Key is not set. This is required for authentication.")
		return
	}
	if data.APIKey.IsNull() {
		apiKey = os.Getenv("POWERDNS_API_KEY")
	} else {
		apiKey = data.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("API Key is not set", "API Key is not set. This is required for authentication.")
		return
	}

	var serverURL string
	if data.ServerURL.IsUnknown() {
		resp.Diagnostics.AddWarning("Server URL is not set", "Server URL is not set. Can't connect to PowerDNS API.")
		return
	}
	if data.ServerURL.IsNull() {
		serverURL = os.Getenv("POWERDNS_SERVER_URL")
	} else {
		serverURL = data.ServerURL.ValueString()
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

	client := powerdns.New(ctx, apiKey, parsedServerURL.Host, parsedServerURL.Path, parsedServerURL.Scheme)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PowerdnsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRecordsetResource,
		NewZoneResource,
	}
}

func (p *PowerdnsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRecordsetDataSource,
		NewZoneDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PowerdnsProvider{
			version: version,
		}
	}
}
