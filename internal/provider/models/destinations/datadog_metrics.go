package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type DatadogMetricsDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	ApiKey       String `tfsdk:"api_key" user_config:"true"`
	Site         String `tfsdk:"site" user_config:"true"`
	AckEnabled   Bool   `tfsdk:"ack_enabled" user_config:"true"`
}

func DatadogMetricsDestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Publishes metric events to Datadog",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Sensitive:   true,
				Required:    true,
				Description: "Datadog metrics application API key.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"site": schema.StringAttribute{
				Required:    true,
				Description: "The Datadog site (region) to send metrics to.",
				Validators:  []validator.String{stringvalidator.OneOf("us1", "us3", "us5", "eu1")},
			},
		}, nil),
	}
}

func DatadogMetricsFromModel(plan *DatadogMetricsDestinationModel, previousState *DatadogMetricsDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Destination{
		BaseNode: BaseNode{
			Type:        "datadog-metrics",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"api_key":     plan.ApiKey.ValueString(),
				"site":        plan.Site.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
			},
		},
	}

	if !plan.Inputs.IsUnknown() {
		component.Inputs = modelutils.StringListValueToStringSlice(plan.Inputs)
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func DatadogMetricsDestinationToModel(plan *DatadogMetricsDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	if component.UserConfig["api_key"] != nil {
		value, _ := component.UserConfig["api_key"].(string)
		plan.ApiKey = StringValue(value)
	}
	if component.UserConfig["site"] != nil {
		value, _ := component.UserConfig["site"].(string)
		plan.Site = StringValue(value)
	}
	if component.UserConfig["ack_enabled"] != nil {
		value, _ := component.UserConfig["ack_enabled"].(bool)
		plan.AckEnabled = BoolValue(value)
	}
}
