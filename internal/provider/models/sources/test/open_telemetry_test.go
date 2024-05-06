package sources

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/providertest"
)

type resourceOpenTelemetrySourceConfig struct {
	// Type is the OpenTelemetry data type: traces, metrics, logs
	Type string

	PipelineName    string
	PipelineTitle   string
	SourceName      string
	SourceTitle     string
	SourceDesc      string
	CaptureMetadata bool

	// Will include a shared source which means setting the gateway_route_id on
	// source A to the ID of source B.
	GatewayResourceID string
}

var resourceOpenTelemetrySourceConfigTpl = `
resource "mezmo_pipeline" "{{.PipelineName}}" {
	title = "{{.PipelineTitle}}"
}

resource "mezmo_open_telemetry_{{.Type}}_source" "{{.SourceName}}" {
	pipeline_id = mezmo_pipeline.{{.PipelineName}}.id
	title       = "{{.SourceTitle}}"
	description = "{{.SourceDesc}}"

	{{if .CaptureMetadata}}
	capture_metadata = "{{.CaptureMetadata}}"
	{{end}}
}

{{if .GatewayResourceID}}
resource "mezmo_open_telemetry_{{.Type}}_source" "{{.SourceName}}-gateway" {
	pipeline_id      = mezmo_pipeline.{{.PipelineName}}.id
	gateway_route_id = {{.GatewayResourceID}}
}
{{end}}
`

func TestAccOpenTelemetrySources(t *testing.T) {
	for _, sourceType := range []string{"traces", "metrics", "logs"} {

		pipelineName := "test_parent"
		pipelineTitle := "parent pipeline"
		sourceName := "my_source"
		sourceTitle := "my open telemetry title"
		sourceDesc := "my open telemetry description"

		resourceName := fmt.Sprintf("mezmo_open_telemetry_%s_source.%s", sourceType, sourceName)
		resourceNamePipeline := fmt.Sprintf("mezmo_pipeline.%s", pipelineName)
		resourceNameGateway := fmt.Sprintf("mezmo_open_telemetry_%s_source.%s-gateway", sourceType, sourceName)

		config, err := ParsedAccConfig(resourceOpenTelemetrySourceConfig{
			Type:          sourceType,
			PipelineName:  pipelineName,
			PipelineTitle: pipelineTitle,
			SourceName:    sourceName,
			SourceTitle:   sourceTitle,
			SourceDesc:    sourceDesc,
		}, resourceOpenTelemetrySourceConfigTpl)
		if err != nil {
			t.Fatalf("error parsing config template: %s", err)
		}

		configUpdate, err := ParsedAccConfig(resourceOpenTelemetrySourceConfig{
			Type:            sourceType,
			PipelineName:    pipelineName,
			PipelineTitle:   pipelineTitle,
			SourceName:      sourceName,
			SourceTitle:     sourceTitle + " updated",
			SourceDesc:      sourceDesc,
			CaptureMetadata: true,
		}, resourceOpenTelemetrySourceConfigTpl)
		if err != nil {
			t.Fatalf("error parsing config template: %s", err)
		}

		configGateway, err := ParsedAccConfig(resourceOpenTelemetrySourceConfig{
			Type:              sourceType,
			PipelineName:      pipelineName,
			PipelineTitle:     pipelineTitle,
			SourceName:        sourceName,
			SourceTitle:       sourceTitle,
			SourceDesc:        sourceDesc,
			GatewayResourceID: resourceName + ".gateway_route_id",
		}, resourceOpenTelemetrySourceConfigTpl)
		if err != nil {
			t.Fatalf("error parsing config template: %s", err)
		}

		configGatewayUpdateErr, err := ParsedAccConfig(resourceOpenTelemetrySourceConfig{
			Type:              sourceType,
			PipelineName:      pipelineName,
			PipelineTitle:     pipelineTitle,
			SourceName:        sourceName,
			SourceTitle:       sourceTitle,
			SourceDesc:        sourceDesc,
			GatewayResourceID: resourceNamePipeline + ".id",
		}, resourceOpenTelemetrySourceConfigTpl)
		if err != nil {
			t.Fatalf("error parsing config template: %s", err)
		}

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			PreCheck:                 func() { TestPreCheck(t) },
			Steps: []resource.TestStep{
				// Error: pipeline_id is required
				{
					Config: fmt.Sprintf(`
					resource "mezmo_pipeline" "test_parent" {
						title = "parent pipeline"
					}
					resource "mezmo_open_telemetry_%s_source" "my_source" {
						title = "my kafka title"
						description = "my kafka description"
					}`, sourceType),
					ExpectError: regexp.MustCompile("The argument \"pipeline_id\" is required"),
				},
				// Create and Read testing
				{
					Config: GetProviderConfig() + config,
					Check: resource.ComposeTestCheckFunc([]resource.TestCheckFunc{
						resource.TestMatchResourceAttr(resourceName, "id", IDRegex),
						resource.TestMatchResourceAttr(resourceName, "gateway_route_id", IDRegex),
						resource.TestCheckResourceAttr(resourceName, "description", "my open telemetry description"),
						resource.TestCheckResourceAttr(resourceName, "generation_id", "0"),
						resource.TestCheckResourceAttr(resourceName, "title", "my open telemetry title"),
						resource.TestCheckResourceAttr(resourceName, "capture_metadata", "false"),
						resource.TestCheckResourceAttrPair(resourceName, "pipeline_id", "mezmo_pipeline.test_parent", "id"),
					}...),
				},
				// Import
				{
					Config:            GetProviderConfig() + config,
					ImportState:       true,
					ResourceName:      resourceName,
					ImportStateIdFunc: ComputeImportId(resourceName),
					ImportStateVerify: true,
				},
				// Update and Read testing
				{
					Config: GetProviderConfig() + configUpdate,
					Check: resource.ComposeTestCheckFunc([]resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "generation_id", "1"),
						resource.TestCheckResourceAttr(resourceName, "title", "my open telemetry title updated"),
						resource.TestCheckResourceAttr(resourceName, "capture_metadata", "true"),
					}...),
				},
				// Supply gateway_route_id
				{
					Config: GetProviderConfig() + configGateway,
					Check: resource.ComposeTestCheckFunc([]resource.TestCheckFunc{
						resource.TestCheckResourceAttrPair(
							resourceNameGateway,
							"gateway_route_id",
							resourceName,
							"gateway_route_id",
						),
					}...),
				},
				// Updating gateway_route_id is not allowed
				{
					Config:      GetProviderConfig() + configGatewayUpdateErr,
					ExpectError: regexp.MustCompile("This field is immutable after resource creation."),
				},
				// confirm manually deleted resources are recreated
				{
					Config: GetProviderConfig() + config,
					Check: resource.ComposeTestCheckFunc(
						TestDeletePipelineNodeManually(
							"mezmo_pipeline.test_parent",
							resourceName,
						),
					),
					// verify resource will be re-created after refresh
					ExpectNonEmptyPlan: true,
				},
			},
		})
	}
}
