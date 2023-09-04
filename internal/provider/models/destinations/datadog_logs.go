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

type DatadogLogsDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	ApiKey       String `tfsdk:"api_key"`
	Site         String `tfsdk:"site"`
	Compression  String `tfsdk:"compression"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled"`
}

func DatadogLogsDestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Publishes log events to Datadog",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Sensitive:   true,
				Required:    true,
				Description: "Datadog logs application API key.",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"site": schema.StringAttribute{
				Required:    true,
				Description: "The Datadog site (region) to send logs to.",
				Validators:  []validator.String{stringvalidator.OneOf("us1", "us3", "us5", "eu1")},
			},
			"compression": schema.StringAttribute{
				Required:    true,
				Description: "The compression strategy used on the encoded data prior to sending..",
				Validators:  []validator.String{stringvalidator.OneOf("none", "gzip")},
			},
		}, nil),
	}
}

func DatadogLogsFromModel(plan *DatadogLogsDestinationModel, previousState *DatadogLogsDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Destination{
		BaseNode: BaseNode{
			Type:        "datadog-logs",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      modelutils.StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"api_key":     plan.ApiKey.ValueString(),
				"site":        plan.Site.ValueString(),
				"compression": plan.Compression.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func DatadogLogsDestinationToModel(plan *DatadogLogsDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	plan.Title = StringValue(component.Title)
	plan.Description = StringValue(component.Description)
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	plan.ApiKey = StringValue(component.UserConfig["api_key"].(string))
	plan.Site = StringValue(component.UserConfig["site"].(string))
	plan.Compression = StringValue(component.UserConfig["compression"].(string))
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
}
