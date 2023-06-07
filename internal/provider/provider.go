package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

var _ provider.Provider = &MezmoProvider{}

// MezmoProvider defines the provider implementation.
type MezmoProvider struct {
	version string
}

// MezmoProviderModel describes the provider data model.
type MezmoProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *MezmoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mezmo"
	resp.Version = p.version
}

func (p *MezmoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Mezmo provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *MezmoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MezmoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	c := client.NewClient()
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *MezmoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPipelineResource,
	}
}

func (p *MezmoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MezmoProvider{
			version: version,
		}
	}
}
