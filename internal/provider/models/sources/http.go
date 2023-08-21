package sources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type HttpSourceModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	Decoding        String `tfsdk:"decoding"`
	CaptureMetadata Bool   `tfsdk:"capture_metadata"`
	GatewayRouteId  String `tfsdk:"gateway_route_id"`
}

func HttpSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents an HTTP source.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"decoding": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("json"),
				Description: "The decoding method for converting frames into data events.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"bytes", "json", "ndjson"),
				},
			},
		}, []string{"capture_metadata", "gateway_route_id"}),
	}
}

func HttpSourceFromModel(plan *HttpSourceModel, previousState *HttpSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        "http",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"decoding":         plan.Decoding.ValueString(),
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

func HttpSourceToModel(plan *HttpSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	if component.UserConfig["format"] != nil {
		decoding, _ := component.UserConfig["decoding"].(string)
		plan.Decoding = StringValue(decoding)
		captureMetadata, _ := component.UserConfig["capture_metadata"].(bool)
		plan.CaptureMetadata = BoolValue(captureMetadata)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.GatewayRouteId = StringValue(component.GatewayRouteId)
}
