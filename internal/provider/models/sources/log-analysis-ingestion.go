package sources

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
)

const LOG_ANALYSIS_INGESTION_SOURCE_TYPE_NAME = "log_analysis_ingestion"
const LOG_ANALYSIS_INGESTION_SOURCE_NODE_NAME = "log-analysis-ingestion"

type LogAnalysisIngestionSourceModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
}

var LogAnalysisIngestionSourceResourceSchema = schema.Schema{
	Description: "Redirect all data sent to the logs.mezmo.com endpoint directly to a pipeline.",
	Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}, nil),
}

func LogAnalysisIngestionSourceFromModel(plan *LogAnalysisIngestionSourceModel, previousState *LogAnalysisIngestionSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Source{
		BaseNode: BaseNode{
			Type:        LOG_ANALYSIS_INGESTION_SOURCE_NODE_NAME,
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

func LogAnalysisIngestionSourceToModel(plan *LogAnalysisIngestionSourceModel, component *Source) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
}
