package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

var _ provider.Provider = &MezmoProvider{}

// MezmoProvider defines the provider implementation.
type MezmoProvider struct {
	version string
}

// MezmoProviderModel describes the provider data model.
type MezmoProviderModel struct {
	Endpoint String `tfsdk:"endpoint"`
	AuthKey  String `tfsdk:"auth_key"`
	Headers  Map    `tfsdk:"headers"`
}

func (p *MezmoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mezmo"
	resp.Version = p.version
}

func (p *MezmoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Mezmo Terraform Provider allows organizations to manage Pipelines" +
			" (sources, processors and destinations) programmatically via Terraform.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "Mezmo API endpoint containing the url scheme, host and port",
				Optional:    true,
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"auth_key": schema.StringAttribute{
				Description: "The authentication key",
				Required:    true,
				Sensitive:   true,
			},
			"headers": schema.MapAttribute{
				Description: "Optional map of headers to send in each request",
				Optional:    true,
				ElementType: StringType,
				Validators: []validator.Map{
					mapvalidator.All(
						mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
						mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
					),
				},
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

	endpoint := "https://api.mezmo.com"
	headers := make(map[string]string)

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}
	if !data.Headers.IsNull() {
		for k, v := range data.Headers.Elements() {
			headers[k] = v.(String).ValueString()
		}
	}

	c := client.NewClient(endpoint, data.AuthKey.ValueString(), headers)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *MezmoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPipelineResource,

		// Sources
		NewAgentSourceResource,
		NewAzureEventHubSourceResource,
		NewDemoSourceResource,
		NewFluentSourceResource,
		NewHttpSourceResource,
		NewKafkaSourceResource,
		NewKinesisFirehoseSourceResource,
		NewLogStashSourceResource,
		NewPrometheusRemoteWriteSourceResource,
		NewS3SourceResource,
		NewSplunkHecSourceResource,
		NewSQSSourceResource,

		// Processors
		NewCompactFieldsProcessorResource,
		NewDecryptFieldsProcessorResource,
		NewDedupeProcessorResource,
		NewDropFieldsProcessorResource,
		NewEncryptFieldsProcessorResource,
		NewFlattenFieldsProcessorResource,
		NewParseProcessorResource,
		NewSampleProcessorResource,
		NewScriptExecutionProcessorResource,
		NewStringifyProcessorResource,
		NewUnrollProcessorResource,

		// Destinations
		NewAzureBlobStorageDestinationResource,
		NewBlackholeDestinationResource,
		NewDatadogLogsDestinationResource,
		NewDatadogMetricsDestinationResource,
		NewElasticSearchDestinationResource,
		NewHoneycombLogsDestinationResource,
		NewHttpDestinationResource,
		NewKafkaDestinationResource,
		NewLokiDestinationResource,
		NewMezmoDestinationResource,
		NewNewRelicDestinationResource,
		NewPrometheusRemoteWriteDestinationResource,
		NewSplunkHecLogsDestinationResource,
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
