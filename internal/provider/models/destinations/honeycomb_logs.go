package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type HoneycombLogsDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled" user_config:"true"`
	DataSet      String `tfsdk:"dataset" user_config:"true"`
	ApiKey       String `tfsdk:"api_key" user_config:"true"`
}

func HoneycombLogsResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Send log data to Honeycomb",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Honeycomb API key",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"dataset": schema.StringAttribute{
				Required:    true,
				Description: "The name of the targeted dataset",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		}, nil),
	}
}

func HoneycombLogsFromModel(plan *HoneycombLogsDestinationModel, previousState *HoneycombLogsDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "honeycomb-logs",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled": plan.AckEnabled.ValueBool(),
				"api_key":     plan.ApiKey.ValueString(),
				"dataset":     plan.DataSet.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func HoneycombLogsToModel(plan *HoneycombLogsDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.ApiKey = StringValue(component.UserConfig["api_key"].(string))
	plan.DataSet = StringValue(component.UserConfig["dataset"].(string))
}
