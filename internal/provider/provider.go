package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mezmo/terraform-provider-mezmo/internal/client"
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
		NewLogAnalysisSourceResource,
		NewLogStashSourceResource,
		NewPrometheusRemoteWriteSourceResource,
		NewS3SourceResource,
		NewSplunkHecSourceResource,
		NewSQSSourceResource,
		NewOpenTelemetryTracesSourceResource,

		// Processors
		NewCompactFieldsProcessorResource,
		NewDecryptFieldsProcessorResource,
		NewDedupeProcessorResource,
		NewDropFieldsProcessorResource,
		NewEncryptFieldsProcessorResource,
		NewFlattenFieldsProcessorResource,
		NewParseProcessorResource,
		NewParseSequentiallyProcessorResource,
		NewRouteProcessorResource,
		NewReduceProcessorResource,
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
		NewGcpCloudStorageDestinationResource,
		NewHoneycombLogsDestinationResource,
		NewHttpDestinationResource,
		NewKafkaDestinationResource,
		NewLokiDestinationResource,
		NewMezmoDestinationResource,
		NewNewRelicDestinationResource,
		NewPrometheusRemoteWriteDestinationResource,
		NewS3DestinationResource,
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

// This function is used only in testing to protect us from forgetting to overwrite plan
// values with what comes back from the API responses. Our models can be complicated, and
// it's possible we will miss things in the provider code and in code review. What's worse is that
// test assertions will actually give false positives even if properties have not been updated.
// @see: https://mezmo.atlassian.net/browse/LOG-18104
func NullifyPlanFields[M ComponentModel](plan *M, schema resourceSchema.Schema) {
	if os.Getenv("TF_ACC") != "1" {
		// Don't do this in prod because modifying the plan like this could have unintended side effects after upgrades
		return
	}

	modelType := reflect.TypeOf(*plan)
	structVal := reflect.ValueOf(plan)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		val, isTagged := field.Tag.Lookup("user_config")
		if isTagged == false || val != "true" {
			// not a user_config field
			continue
		}
		schemaFieldName := field.Tag.Get("tfsdk") // guaranteed to exist
		structFieldName := field.Name
		// Get the TF type from the schema so it can be used for object/list ElemType etc.
		schemaFieldType, _ := schema.TypeAtPath(context.Background(), path.Empty().AtName(schemaFieldName))
		// Get a reference to the struct field so it can be overwritten
		structField := structVal.Elem().FieldByName(structFieldName)

		var nullValue reflect.Value

		switch schemaFieldType.(type) {
		case basetypes.StringType:
			nullValue = reflect.ValueOf(
				basetypes.NewStringNull(),
			)
		case basetypes.BoolType:
			nullValue = reflect.ValueOf(
				basetypes.NewBoolNull(),
			)
		case basetypes.Int64Type:
			nullValue = reflect.ValueOf(
				basetypes.NewInt64Null(),
			)
		case basetypes.ListType:
			nullValue = reflect.ValueOf(
				basetypes.NewListNull(schemaFieldType.(ListType).ElemType),
			)
		case basetypes.ObjectType:
			nullValue = reflect.ValueOf(
				basetypes.NewObjectNull(schemaFieldType.(ObjectType).AttrTypes),
			)
		case basetypes.MapType:
			nullValue = reflect.ValueOf(
				basetypes.NewMapNull(schemaFieldType.(MapType).ElemType),
			)
		default:
			panic(fmt.Errorf("Unsupported NullifyPlanFields type \"%T\" for field: \"%s\" in schema \"%s\"", schemaFieldType, schemaFieldName, schema.Description))
		}

		structField.Set(nullValue)
	}
}
