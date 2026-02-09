package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/modelutils"
)

const AGGREGATE_PROCESSOR_NODE_NAME = "aggregate"
const AGGREGATE_PROCESSOR_TYPE_NAME = "aggregate"

type AggregateProcessorModel struct {
	Id             String              `tfsdk:"id"`
	PipelineId     String              `tfsdk:"pipeline_id"`
	Title          String              `tfsdk:"title"`
	Description    String              `tfsdk:"description"`
	Inputs         List                `tfsdk:"inputs"`
	GenerationId   Int64               `tfsdk:"generation_id"`
	Interval       Int64               `tfsdk:"interval" user_config:"true"`
	Conditional    Object              `tfsdk:"conditional" user_config:"true"`
	GroupBy        basetypes.ListValue `tfsdk:"group_by" user_config:"true"`
	Script         String              `tfsdk:"script" user_config:"true"`
	WindowMin      Int64               `tfsdk:"window_min" user_config:"true"`
	WindowType     String              `tfsdk:"window_type" user_config:"true"`
	Operation      String              `tfsdk:"operation" user_config:"true"`
	EventTimestamp String              `tfsdk:"event_timestamp" user_config:"true"`
}

var AggregateProcessorResourceSchema = schema.Schema{
	Description: "Aggregates multiple metric events into a single metric event using either a tumbling interval window or a sliding interval window",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"window_type": schema.StringAttribute{
			Required:    true,
			Description: "The type of window to use when aggregating events (tumbling or sliding)",
			Validators:  []validator.String{stringvalidator.OneOf("tumbling", "sliding")},
		},
		"interval": schema.Int64Attribute{
			Required:    true,
			Description: "The interval over which events are aggregated in seconds",
		},
		"operation": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The operation in which to perform the aggregation",
			Validators:  []validator.String{stringvalidator.OneOf(MapKeys(Aggregate_Operations)...)},
		},
		"script": schema.StringAttribute{
			Optional: true,
			Computed: false,
		},
		"window_min": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The minimum amount of time to wait before creating a new sliding window for future events.",
		},
		"conditional": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "When method is set to sliding: " + ParentConditionalAttribute(Non_Change_Operator_Labels).Description,
			Attributes:  ParentConditionalAttribute(Non_Change_Operator_Labels).Attributes,
		},
		"group_by": schema.ListAttribute{
			ElementType: basetypes.StringType{},
			Optional:    true,
			Description: "Group events based on matching data from each of these field paths. Supports nesting via dot-notation.",
		},
		"event_timestamp": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Force accumulated event reduction to flush the result when a conditional expression evaluates to true on an inbound event.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(200),
			},
		},
	}),
}

func windowConfigFromModel(plan *AggregateProcessorModel, userConfig map[string]any, dd *diag.Diagnostics) {
	windowConfig := make(map[string]any)
	windowConfig["type"] = plan.WindowType.ValueString()
	windowConfig["interval"] = plan.Interval.ValueInt64()
	if windowConfig["type"] == "sliding" {
		windowConfig["window_min"] = plan.WindowMin.ValueInt64()
	}
	userConfig["window"] = windowConfig
}

func evaluateConfigFromModel(plan *AggregateProcessorModel, userConfig map[string]any, dd *diag.Diagnostics) {
	isOperationSet := !(plan.Operation.IsNull() || plan.Operation.IsUnknown())
	isScripSet := !(plan.Script.IsNull() || plan.Script.IsUnknown())
	evaluateConfig := make(map[string]any)
	if !isOperationSet && !isScripSet {
		dd.AddError(
			"Error in plan",
			"Either 'operation' or 'script' must be defined.",
		)
	} else if isOperationSet && isScripSet {
		dd.AddError(
			"Error in plan",
			"Cannot define both 'operation' and 'script' fields.",
		)
	} else if isScripSet {
		evaluateConfig["operation"] = "CUSTOM"
		evaluateConfig["script"] = plan.Script.ValueString()
	} else {
		delete(evaluateConfig, "script")
		evaluateConfig["operation"] = Aggregate_Operations[plan.Operation.ValueString()]
	}
	userConfig["evaluate"] = evaluateConfig
}

func AggregateProcessorFromModel(plan *AggregateProcessorModel, previousState *AggregateProcessorModel) (*Processor, diag.Diagnostics) {
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

	windowConfigFromModel(plan, user_config, &dd)
	evaluateConfigFromModel(plan, user_config, &dd)

	if !plan.Conditional.IsNull() {
		user_config["conditional"] = UnwindConditionalFromModel(plan.Conditional)
	}

	if !plan.GroupBy.IsNull() && len(plan.GroupBy.Elements()) > 0 {
		component.UserConfig["group_by"] = StringListValueToStringSlice(plan.GroupBy)
	}

	if !plan.EventTimestamp.IsNull() && !plan.EventTimestamp.IsUnknown() {
		user_config["event_timestamp"] = plan.EventTimestamp.ValueString()
	}

	return &component, dd
}

func AggregateProcessorToModel(plan *AggregateProcessorModel, component *Processor) {
	plan.Id = basetypes.NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = basetypes.NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = basetypes.NewStringValue(component.Description)
	}
	plan.GenerationId = basetypes.NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	windowConfig := component.UserConfig["window"].(map[string]any)
	plan.WindowType = basetypes.NewStringValue(windowConfig["type"].(string))
	plan.Interval = Int64Value(int64(windowConfig["interval"].(float64)))
	if windowConfig["window_min"] != nil {
		plan.WindowMin = Int64Value(int64(windowConfig["window_min"].(float64)))
	}

	evaluateConfig := component.UserConfig["evaluate"].(map[string]any)
	if evaluateConfig["operation"] != nil {
		apiOperation := evaluateConfig["operation"].(string)
		if apiOperation == "CUSTOM" {
			plan.Operation = basetypes.NewStringNull()
			plan.Script = basetypes.NewStringValue(evaluateConfig["script"].(string))
		} else {
			plan.Operation = basetypes.NewStringValue(FindKey(Aggregate_Operations, apiOperation))
			plan.Script = basetypes.NewStringNull()
		}
	}

	if component.UserConfig["conditional"] != nil {
		plan.Conditional = UnwindConditionalToModel(component.UserConfig["conditional"].(map[string]any), Non_Change_Operator_Labels)
	}

	if component.UserConfig["group_by"] != nil {
		plan.GroupBy = SliceToStringListValue(component.UserConfig["group_by"].([]any))
	}

	if component.UserConfig["event_timestamp"] != nil {
		plan.EventTimestamp = basetypes.NewStringValue(component.UserConfig["event_timestamp"].(string))
	}

}
