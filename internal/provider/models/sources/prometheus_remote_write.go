package sources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
)

type PrometheusRemoteWriteSourceModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	GatewayRouteId  String `tfsdk:"gateway_route_id"`
	CaptureMetadata Bool   `tfsdk:"capture_metadata" user_config:"true"`
}

var PrometheusRemoteWriteSourceResourceSchema = schema.Schema{
	Description: "Represents a Prometheus Remote Write source.",
	Attributes:  ExtendBaseAttributes(map[string]schema.Attribute{}, []string{"capture_metadata", "gateway_route_id"}),
}

func PrometheusRemoteWriteSourceFromModel(plan *PrometheusRemoteWriteSourceModel, previousState *PrometheusRemoteWriteSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        "prometheus-remote-write",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"capture_metadata": plan.CaptureMetadata.ValueBool(),
			},
		},
	}

	if previousState == nil {
		if !plan.GatewayRouteId.IsUnknown() {
			// Let them specify gateway route id on POST only
			component.GatewayRouteId = plan.GatewayRouteId.ValueString()
		}
	} else {
		// Set generated fields
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()

		// If they have specified gateway_route_id, then it *cannot* be a different value that what's in state
		if !plan.GatewayRouteId.IsUnknown() && plan.GatewayRouteId.ValueString() != previousState.GatewayRouteId.ValueString() {
			details := fmt.Sprintf(
				"Cannot update \"gateway_route_id\" to %s. This field is immutable after resource creation.",
				plan.GatewayRouteId,
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
	if component.UserConfig["capture_metadata"] != nil {
		plan.CaptureMetadata = BoolValue(component.UserConfig["capture_metadata"].(bool))
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.GatewayRouteId = StringValue(component.GatewayRouteId)
}
