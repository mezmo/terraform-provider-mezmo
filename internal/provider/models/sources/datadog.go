package sources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
)

const DATADOG_SOURCE_NODE_NAME = "mezmo-datadog-source"
const DATADOG_SOURCE_TYPE_NAME = "datadog"

type DatadogSourceModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	SharedSourceId  String `tfsdk:"shared_source_id"`
	CaptureMetadata Bool   `tfsdk:"capture_metadata" user_config:"true"`
}

var DatadogSourceResourceSchema = schema.Schema{
	Description: "Send logs and metrics data directly from an installed datadog agent",
	Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}, []string{"capture_metadata", "shared_source_id"}),
}

func DatadogSourceFromModel(plan *DatadogSourceModel, previousState *DatadogSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        DATADOG_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"capture_metadata": plan.CaptureMetadata.ValueBool(),
			},
		},
	}

	if previousState == nil {
		if !plan.SharedSourceId.IsUnknown() {
			// Let them specify gateway route id on POST only
			component.SharedSourceId = plan.SharedSourceId.ValueString()
		}
	} else {
		// Set generated fields
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()

		// If they have specified shared_source_id, then it *cannot* be a different value that what's in state
		if !plan.SharedSourceId.IsUnknown() && plan.SharedSourceId.ValueString() != previousState.SharedSourceId.ValueString() {
			details := fmt.Sprintf(
				"Cannot update \"shared_source_id\" to %s. This field is immutable after resource creation.",
				plan.SharedSourceId,
			)
			dd.AddError("Error in plan", details)
			return nil, dd
		}
	}

	return &component, dd
}

func DatadogSourceToModel(plan *DatadogSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.CaptureMetadata = BoolValue(component.UserConfig["capture_metadata"].(bool))
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.SharedSourceId = StringValue(component.SharedSourceId)
}
