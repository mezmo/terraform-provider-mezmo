package sources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
)

const PROMETHEUS_REMOTE_WRITE_SOURCE_TYPE_NAME = "prometheus_remote_write"
const PROMETHEUS_REMOTE_WRITE_SOURCE_NODE_NAME = "prometheus-remote-write"

type PrometheusRemoteWriteSourceModel struct {
	Id             String `tfsdk:"id"`
	PipelineId     String `tfsdk:"pipeline_id"`
	Title          String `tfsdk:"title"`
	Description    String `tfsdk:"description"`
	GenerationId   Int64  `tfsdk:"generation_id"`
	SharedSourceId String `tfsdk:"shared_source_id"`
}

var PrometheusRemoteWriteSourceResourceSchema = schema.Schema{
	Description: "Represents a Prometheus Remote Write source.",
	Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}, []string{"shared_source_id"}),
}

func PrometheusRemoteWriteSourceFromModel(plan *PrometheusRemoteWriteSourceModel, previousState *PrometheusRemoteWriteSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        PROMETHEUS_REMOTE_WRITE_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{},
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

func PrometheusRemoteWriteSourceToModel(plan *PrometheusRemoteWriteSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.SharedSourceId = StringValue(component.SharedSourceId)
}
