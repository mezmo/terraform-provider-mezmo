package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type DedupeProcessorModel struct {
	Id             String `tfsdk:"id"`
	PipelineId     String `tfsdk:"pipeline_id"`
	Title          String `tfsdk:"title"`
	Description    String `tfsdk:"description"`
	Inputs         List   `tfsdk:"inputs"`
	GenerationId   Int64  `tfsdk:"generation_id"`
	Fields         List   `tfsdk:"fields" user_config:"true"`
	NumberOfEvents Int64  `tfsdk:"number_of_events" user_config:"true"`
	ComparisonType String `tfsdk:"comparison_type" user_config:"true"`
}

func DedupeProcessorResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Remove duplicates from the data stream",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"fields": schema.ListAttribute{
				ElementType: StringType,
				Required:    true,
				Description: "A list of fields on which to base deduping",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"number_of_events": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: "Number of events to compare across",
				Validators: []validator.Int64{
					int64validator.AtLeast(2),
					int64validator.AtMost(5000),
				},
				Default: int64default.StaticInt64(5000),
			},
			"comparison_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Description: "When set to \"Match\" (default), it only compares across the fields which are" +
					" specified by the user. When set to \"Ignore\", it compares everything but the fields" +
					" specified by the user",
				Default: stringdefault.StaticString("Match"),
				Validators: []validator.String{
					stringvalidator.OneOf("Ignore", "Match"),
				},
			},
		}),
	}
}

func DedupeProcessorFromModel(plan *DedupeProcessorModel, previousState *DedupeProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "dedupe",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  make(map[string]any),
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)
	component.UserConfig["fields"] = StringListValueToStringSlice(plan.Fields)
	// Default values make the plan always to have this values defined
	component.UserConfig["number_of_events"] = plan.NumberOfEvents.ValueInt64()
	component.UserConfig["comparison_type"] = plan.ComparisonType.ValueString()

	return &component, dd
}

func DedupeProcessorToModel(plan *DedupeProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Fields = SliceToStringListValue(component.UserConfig["fields"].([]any))
	if component.UserConfig["number_of_events"] != nil {
		plan.NumberOfEvents = Int64Value(int64(component.UserConfig["number_of_events"].(float64)))
	}
	if component.UserConfig["comparison_type"] != nil {
		plan.ComparisonType = StringValue(component.UserConfig["comparison_type"].(string))
	}
}
