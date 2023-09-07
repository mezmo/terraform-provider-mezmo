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

type KinesisFirehoseSourceModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	GatewayRouteId  String `tfsdk:"gateway_route_id"`
	Decoding        String `tfsdk:"decoding" user_config:"true"`
	CaptureMetadata Bool   `tfsdk:"capture_metadata" user_config:"true"`
}

func KinesisFirehoseSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Receive Kinesis Firehose data",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"decoding": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("json"),
				Description: "This specifies what the data format will be after it is base64 decoded. " +
					"If it is JSON, it will be automatically parsed.",
				Validators: []validator.String{
					stringvalidator.OneOf("text", "json"),
				},
			},
		}, []string{"capture_metadata", "gateway_route_id"}),
	}
}

func KinesisFirehoseSourceFromModel(plan *KinesisFirehoseSourceModel, previousState *KinesisFirehoseSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Source{
		BaseNode: BaseNode{
			Type:        "kinesis-firehose",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"format":           plan.Decoding.ValueString(),
				"capture_metadata": plan.CaptureMetadata.ValueBool(),
			},
		},
	}

	if previousState == nil {
		if !plan.GatewayRouteId.IsUnknown() {
			component.GatewayRouteId = plan.GatewayRouteId.ValueString()
		}
	} else {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()

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

func KinesisFirehoseSourceToModel(plan *KinesisFirehoseSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.Decoding = StringValue(component.UserConfig["format"].(string))
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.CaptureMetadata = BoolValue(component.UserConfig["capture_metadata"].(bool))
	plan.GatewayRouteId = StringValue(component.GatewayRouteId)
}
