package sources

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
)

type LogAnalysisSourceModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
}

func LogAnalysisSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Receive data directly from your Mezmo Log Analysis account",
		Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}, nil),
	}
}

func LogAnalysisSourceFromModel(plan *LogAnalysisSourceModel, previousState *LogAnalysisSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Source{
		BaseNode: BaseNode{
			Type:        "log-analysis",
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

func LogAnalysisSourceToModel(plan *LogAnalysisSourceModel, component *Source) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
}
