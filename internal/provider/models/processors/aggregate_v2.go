package processors

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const AGGREGATE_PROCESSOR_NODE_NAME = "aggregate-v2"
const AGGREGATE_PROCESSOR_TYPE_NAME = "aggregate_v2"

var METHODS = map[string]string{"TUMBLING": "tumbling", "SLIDING": "sliding"}
var STRATEGIES = map[string]string{"SUM": "SUM", "AVERAGE": "AVG", "SET_INTERSECTION": "SET_INTERSECTION", "DISTROBUTION_CONCATENATION": "DIST_CONCAT"}

type AggregateV2ProcessorModel struct {
	Id           String              `tfsdk:"id"`
	PipelineId   String              `tfsdk:"pipeline_id"`
	Title        String              `tfsdk:"title"`
	Description  String              `tfsdk:"description"`
	Inputs       List                `tfsdk:"inputs"`
	GenerationId Int64               `tfsdk:"generation_id"`
	Method       String              `tfsdk:"method" user_config:"true"`
	IntervalMS   Int64               `tfsdk:"interval" user_config:"true"`
	Strategy     String              `tfsdk:"strategy" user_config:"true"`
	Duration     Int64               `tfsdk:"window_duration" user_config:"true"`
	Conditional  Object              `tfsdk:"conditional" user_config:"true"`
	GroupBy      basetypes.ListValue `tfsdk:"group_by" user_config:"true"`
}

var AggregateV2ProcessorResourceSchema = schema.Schema{
	Description: "Aggregates multiple metric events into a single metric event using either a tumbling interval window or a sliding interval window",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"method": schema.StringAttribute{
			Required:    true,
			Description: "The method in which to aggregate metrics (tumbling or sliding)",
		},
		"interval": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "When method is set to tumbling, this is the interval over which metrics are aggregated in seconds",
		},
		"strategy": schema.StringAttribute{
			Required:    true,
			Description: "When method is set to sliding, this is the strategy in which to perform the aggregation",
		},
		"window_duration": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "When method is set to sliding, this is the interval over which metrics are aggregated in seconds",
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

	if !plan.Method.IsNull() {
		user_config["method"] = plan.Method.ValueString()
	}

	if !plan.IntervalMS.IsNull() {
		user_config["interval"] = plan.IntervalMS.ValueInt64()
	}

	if !plan.Strategy.IsNull() {
		user_config["strategy"] = plan.Strategy.ValueString()
	}

	if !plan.Duration.IsNull() {
		user_config["window_duration"] = plan.Duration.ValueInt64()
	}

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
		plan.IntervalMS = Int64Value(int64(component.UserConfig["interval"].(float64)))
	}

	if component.UserConfig["strategy"] != nil {
		plan.Strategy = basetypes.NewStringValue(component.UserConfig["strategy"].(string))
	}

	if component.UserConfig["window_duration"] != nil {
		plan.Duration = Int64Value(int64(component.UserConfig["window_duration"].(float64)))
	}

	if component.UserConfig["conditional"] != nil {
		plan.Conditional = UnwindConditionalToModel(component.UserConfig["conditional"].(map[string]any))
	}

	if component.UserConfig["group_by"] != nil {
		plan.GroupBy = SliceToStringListValue(component.UserConfig["group_by"].([]any))
	}

}
