package sources

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v6/internal/client"
)

const OPEN_TELEMETRY_COMBINED_SOURCE_TYPE_NAME = "open_telemetry_combined"
const OPEN_TELEMETRY_COMBINED_SOURCE_NODE_NAME = "open-telemetry"

type OpenTelemetryCombinedSourceModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
}

var OpenTelemetryCombinedSourceResourceSchema = schema.Schema{
	Description: "Receive logs, metrics & traces from OpenTelemetry senders, authenticated with your " +
		"account's ingestion key.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{}, nil),
}

func OpenTelemetryCombinedSourceFromModel(plan *OpenTelemetryCombinedSourceModel, previousState *OpenTelemetryCombinedSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Source{
		BaseNode: BaseNode{
			Type:        OPEN_TELEMETRY_COMBINED_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func OpenTelemetryCombinedSourceToModel(plan *OpenTelemetryCombinedSourceModel, component *Source) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
}
