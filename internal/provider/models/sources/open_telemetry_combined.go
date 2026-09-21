package sources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v6/internal/client"
)

const OPEN_TELEMETRY_COMBINED_SOURCE_TYPE_NAME = "open_telemetry_combined"
const OPEN_TELEMETRY_COMBINED_SOURCE_NODE_NAME = "open-telemetry"

type OpenTelemetryCombinedSourceModel struct {
	Id              StringValue `tfsdk:"id"`
	PipelineId      StringValue `tfsdk:"pipeline_id"`
	Title           StringValue `tfsdk:"title"`
	Description     StringValue `tfsdk:"description"`
	GenerationId    Int64Value  `tfsdk:"generation_id"`
	LogsOutputId    StringValue `tfsdk:"logs_output_id"`
	MetricsOutputId StringValue `tfsdk:"metrics_output_id"`
	TracesOutputId  StringValue `tfsdk:"traces_output_id"`
}

var OpenTelemetryCombinedSourceResourceSchema = schema.Schema{
	Description: "Receive logs, metrics & traces from OpenTelemetry senders, authenticated with your " +
		"account's ingestion key.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"logs_output_id": schema.StringAttribute{
			Computed: true,
			Description: "A system generated value to identify the logs produced by this source. " +
				"This value should be used when connecting the logs to a processor or destination.",
		},
		"metrics_output_id": schema.StringAttribute{
			Computed: true,
			Description: "A system generated value to identify the metrics produced by this source. " +
				"This value should be used when connecting the metrics to a processor or destination.",
		},
		"traces_output_id": schema.StringAttribute{
			Computed: true,
			Description: "A system generated value to identify the traces produced by this source. " +
				"This value should be used when connecting the traces to a processor or destination.",
		},
	}, nil),
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

	outputs := mappedOtelCombinedOutputs(component)
	if id, ok := outputs["logs"]; ok {
		plan.LogsOutputId = id
	}
	if id, ok := outputs["metrics"]; ok {
		plan.MetricsOutputId = id
	}
	if id, ok := outputs["traces"]; ok {
		plan.TracesOutputId = id
	}
}

// mappedOtelCombinedOutputs indexes a source's named outputs (each returned by the
// API as "<component_id>.<output_name>") by their lower-cased output name.
func mappedOtelCombinedOutputs(component *Source) map[string]StringValue {
	res := make(map[string]StringValue, len(component.Outputs))
	for _, item := range component.Outputs {
		idParts := strings.Split(item.Id, ".")
		if len(idParts) == 2 {
			res[strings.ToLower(idParts[1])] = NewStringValue(item.Id)
		}
	}
	return res
}
