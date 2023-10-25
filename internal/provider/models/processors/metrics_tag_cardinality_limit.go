package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type MetricsTagCardinalityLimitProcessorModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	Inputs       ListValue   `tfsdk:"inputs"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
	Tags         ListValue   `tfsdk:"tags" user_config:"true"`
	ExcludeTags  ListValue   `tfsdk:"exclude_tags" user_config:"true"`
	Action       StringValue `tfsdk:"action" user_config:"true"`
	ValueLimit   Int64Value  `tfsdk:"value_limit" user_config:"true"`
	Mode         StringValue `tfsdk:"mode" user_config:"true"`
}

var MetricsTagCardinalityLimitProcessorName = "metrics_tag_cardinality_limit"

var MetricsTagCardinalityLimitProcessorResourceSchema = schema.Schema{
	Description: "Limits the cardinality of metric events by either dropping events " +
		"or tags that exceed a specified value limit",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"tags": schema.ListAttribute{
			ElementType: StringType{},
			Optional:    true,
			Description: "A list of tags to apply cardinality limits. If none are provided, " +
				"all tags will be considered.",
			Validators: []validator.List{
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtMost(100)),
				listvalidator.SizeAtMost(10),
			},
		},
		"exclude_tags": schema.ListAttribute{
			ElementType: StringType{},
			Optional:    true,
			Description: "A list of tags to explicitly exclude from cardinality limits",
			Validators: []validator.List{
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtMost(100)),
				listvalidator.SizeAtMost(10),
			},
		},
		"action": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("drop_event"),
			Description: "The action to take when a tag's cardinality exceeds the value limit",
			Validators: []validator.String{
				stringvalidator.OneOf(LimitExceedAction...),
			},
		},
		"value_limit": schema.Int64Attribute{
			Required:    true,
			Description: "Maximum number of unique values for tags",
		},
		"mode": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("exact"),
			Description: "The method to used to reduce tag value cardinality",
			Validators: []validator.String{
				stringvalidator.OneOf(TagCardinalityMode...),
			},
		},
	}),
}

func MetricsTagCardinalityLimitProcessorFromModel(plan *MetricsTagCardinalityLimitProcessorModel, previousState *MetricsTagCardinalityLimitProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := &Processor{
		BaseNode: BaseNode{
			Type:        "metrics-tag-cardinality-limit",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	user_config := component.UserConfig

	if !plan.Tags.IsNull() {
		user_config["tags"] = StringListValueToStringSlice(plan.Tags)
	}
	if !plan.ExcludeTags.IsNull() {
		user_config["exclude_tags"] = StringListValueToStringSlice(plan.ExcludeTags)
	}
	if !plan.Action.IsNull() {
		user_config["action"] = plan.Action.ValueString()
	}
	if !plan.ValueLimit.IsNull() {
		user_config["value_limit"] = plan.ValueLimit.ValueInt64()
	}
	if !plan.Mode.IsNull() {
		user_config["mode"] = plan.Mode.ValueString()
	}

	return component, dd
}

func MetricsTagCardinalityLimitProcessorToModel(plan *MetricsTagCardinalityLimitProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	if component.UserConfig["tags"] != nil {
		plan.Tags = SliceToStringListValue(component.UserConfig["tags"].([]any))
	}
	if component.UserConfig["exclude_tags"] != nil {
		plan.ExcludeTags = SliceToStringListValue(component.UserConfig["exclude_tags"].([]any))
	}
	if component.UserConfig["action"] != nil {
		plan.Action = NewStringValue(component.UserConfig["action"].(string))
	}
	if component.UserConfig["value_limit"] != nil {
		plan.ValueLimit = NewInt64Value(int64(component.UserConfig["value_limit"].(float64)))
	}
	if component.UserConfig["mode"] != nil {
		plan.Mode = NewStringValue(component.UserConfig["mode"].(string))
	}
}
