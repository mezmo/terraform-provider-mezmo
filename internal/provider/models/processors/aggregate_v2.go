package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const AGGREGATE_PROCESSOR_NODE_NAME = "aggregate-v2"
const AGGREGATE_PROCESSOR_TYPE_NAME = "aggregate_v2"

var STRATEGIES = map[string]string{
	"sum":                        "SUM",
	"average":                    "AVG",
	"set_intersection":           "SET_INTERSECTION",
	"distribution_concatenation": "DIST_CONCAT",
}

type AggregateV2ProcessorModel struct {
	Id           String              `tfsdk:"id"`
	PipelineId   String              `tfsdk:"pipeline_id"`
	Title        String              `tfsdk:"title"`
	Description  String              `tfsdk:"description"`
	Inputs       List                `tfsdk:"inputs"`
	GenerationId Int64               `tfsdk:"generation_id"`
	Method       String              `tfsdk:"method" user_config:"true"`
	Interval     Int64               `tfsdk:"interval" user_config:"true"`
	Strategy     String              `tfsdk:"strategy" user_config:"true"`
	Duration     Int64               `tfsdk:"window_duration" user_config:"true"`
	Minimum      Int64               `tfsdk:"window_min" user_config:"true"`
	Conditional  Object              `tfsdk:"conditional" user_config:"true"`
	GroupBy      basetypes.ListValue `tfsdk:"group_by" user_config:"true"`
	Script       String              `tfsdk:"script" user_config:"true"`
}

var AggregateV2ProcessorResourceSchema = schema.Schema{
	Description: "Aggregates multiple metric events into a single metric event using either a tumbling interval window or a sliding interval window",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"method": schema.StringAttribute{
			Required:    true,
			Description: "The method in which to aggregate metrics (tumbling or sliding)",
			Validators:  []validator.String{stringvalidator.OneOf("tumbling", "sliding")},
		},
		"interval": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "When method is set to tumbling, this is the interval over which metrics are aggregated in seconds",
		},
		"strategy": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "When method is set to sliding, this is the strategy in which to perform the aggregation",
			Validators:  []validator.String{stringvalidator.OneOf(MapKeys(STRATEGIES)...)},
		},
		"script": schema.StringAttribute{
			Optional: true,
			Computed: false,
		},
		"window_duration": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "When method is set to sliding, this is the interval over which metrics are aggregated in seconds",
		},
		"window_min": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "",
		},
		"conditional": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "When method is set to sliding: " + ParentConditionalAttribute.Description,
			Attributes:  ParentConditionalAttribute.Attributes,
		},
		"group_by": schema.ListAttribute{
			ElementType: basetypes.StringType{},
			Optional:    true,
			Computed:    true,
			Description: "Group events based on matching data from each of these field paths. Supports nesting via dot-notation.",
		},
	}),
}

func methodConfigFromModel(plan *AggregateV2ProcessorModel, userConfig map[string]any, dd *diag.Diagnostics) {
	isIntervalSet := !plan.Interval.IsNull() && !plan.Interval.IsUnknown()
	isDurationSet := !plan.Duration.IsNull() && !plan.Duration.IsUnknown()
	isMinimumSet := !plan.Minimum.IsNull() && !plan.Minimum.IsUnknown()
	method := plan.Method.ValueString()

	userConfig["method"] = method
	if method == "tumbling" {
		if isDurationSet {
			dd.AddError(
				"Error in plan",
				"The field 'window_duration' can only be set if method == 'sliding'",
			)
		}

		if isMinimumSet {
			dd.AddError(
				"Error in plan",
				"The field 'window_min' can only be set if method == 'sliding'",
			)
		}

		if isIntervalSet {
			userConfig["interval"] = plan.Interval.ValueInt64()
		}
	} else if method == "sliding" {
		if isIntervalSet {
			dd.AddError(
				"Error in plan",
				"The field 'interval' can only be set if method == 'tumbling'",
			)
		}

		if isDurationSet {
			userConfig["window_duration"] = plan.Duration.ValueInt64()
		}
		if isMinimumSet {
			userConfig["window_min"] = plan.Minimum.ValueInt64()
		}
	} else {
		dd.AddError(
			"Error in plan",
			"The method '%s' is not handled correctly in the provider. Please open a GitHub issue to report this.",
		)
	}
}

func strategyConfigFromModel(plan *AggregateV2ProcessorModel, userConfig map[string]any, dd *diag.Diagnostics) {
	strategySet := !(plan.Strategy.IsNull() || plan.Strategy.IsUnknown())
	scripSet := !(plan.Script.IsNull() || plan.Script.IsNull())
	if !strategySet && !scripSet {
		dd.AddError(
			"Error in plan",
			"Either 'strategy' or 'script' must be defined.",
		)
	} else if strategySet && scripSet {
		dd.AddError(
			"Error in plan",
			"Cannot define both 'strategy' and 'script' fields.",
		)
	} else if scripSet {
		userConfig["strategy"] = "CUSTOM"
		userConfig["script"] = plan.Script.ValueString()
	} else {
		delete(userConfig, "script")
		userConfig["strategy"] = STRATEGIES[plan.Strategy.ValueString()]
	}

}

func AggregateV2ProcessorFromModel(plan *AggregateV2ProcessorModel, previousState *AggregateV2ProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        AGGREGATE_PROCESSOR_NODE_NAME,
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
	user_config := component.UserConfig

	methodConfigFromModel(plan, user_config, &dd)
	strategyConfigFromModel(plan, user_config, &dd)

	if !plan.Conditional.IsNull() {
		user_config["conditional"] = unwindConditionalFromModel(plan.Conditional)
	}

	if !plan.GroupBy.IsNull() && len(plan.GroupBy.Elements()) > 0 {
		component.UserConfig["group_by"] = StringListValueToStringSlice(plan.GroupBy)
	}

	return &component, dd
}

func AggregateV2ProcessorToModel(plan *AggregateV2ProcessorModel, component *Processor) {
	plan.Id = basetypes.NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = basetypes.NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = basetypes.NewStringValue(component.Description)
	}
	plan.GenerationId = basetypes.NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Method = basetypes.NewStringValue(component.UserConfig["method"].(string))

	if component.UserConfig["interval"] != nil {
		plan.Interval = Int64Value(int64(component.UserConfig["interval"].(float64)))
	}

	if component.UserConfig["strategy"] != nil {
		apiStrategy := component.UserConfig["strategy"].(string)
		if apiStrategy == "CUSTOM" {
			plan.Strategy = basetypes.NewStringNull()
			plan.Script = basetypes.NewStringValue(component.UserConfig["script"].(string))
		} else {
			plan.Strategy = basetypes.NewStringValue(FindKey(STRATEGIES, apiStrategy))
			plan.Script = basetypes.NewStringNull()
		}
	}

	if component.UserConfig["window_duration"] != nil {
		plan.Duration = Int64Value(int64(component.UserConfig["window_duration"].(float64)))
	}

	if component.UserConfig["window_min"] != nil {
		plan.Minimum = Int64Value(int64(component.UserConfig["window_min"].(float64)))
	}

	if component.UserConfig["conditional"] != nil {
		plan.Conditional = UnwindConditionalToModel(component.UserConfig["conditional"].(map[string]any))
	}

	if component.UserConfig["group_by"] != nil {
		plan.GroupBy = SliceToStringListValue(component.UserConfig["group_by"].([]any))
	}

}
