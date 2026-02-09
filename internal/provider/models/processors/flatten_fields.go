package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/modelutils"
)

const FLATTEN_FIELDS_PROCESSOR_NODE_NAME = "flatten-fields"
const FLATTEN_FIELDS_PROCESSOR_TYPE_NAME = "flatten_fields"

type FlattenFieldsProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Fields       List   `tfsdk:"fields" user_config:"true"`
	Delimiter    String `tfsdk:"delimiter" user_config:"true"`
}

var FlattenFieldsProcessorResourceSchema = schema.Schema{
	Description: "Flattens the object or array value of a field into a single-level representation.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"fields": schema.ListAttribute{
			ElementType: StringType,
			Optional:    true,
			Description: "A list of nested fields containing a value to flatten. When empty or omitted, the entire event will be flattened.",
			Validators: []validator.List{
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"delimiter": schema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The separator to use between flattened field names",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(1),
			},
			Default: stringdefault.StaticString("_"),
		},
	}),
}

func FlattenFieldsProcessorFromModel(plan *FlattenFieldsProcessorModel, previousState *FlattenFieldsProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        FLATTEN_FIELDS_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  make(map[string]any),
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	if !plan.Inputs.IsUnknown() {
		component.Inputs = modelutils.StringListValueToStringSlice(plan.Inputs)
	}

	component.UserConfig["fields"] = modelutils.StringListValueToStringSlice(plan.Fields)
	component.UserConfig["options"] = make(map[string]any)
	component.UserConfig["options"].(map[string]any)["delimiter"] = plan.Delimiter.ValueString()

	return &component, dd
}

func FlattenFieldsProcessorToModel(plan *FlattenFieldsProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	if component.Inputs != nil {
		inputs := make([]attr.Value, 0)
		for _, v := range component.Inputs {
			inputs = append(inputs, StringValue(v))
		}
		plan.Inputs = ListValueMust(StringType, inputs)
	}

	if component.UserConfig["fields"] != nil {
		fields, _ := component.UserConfig["fields"].([]any)
		if len(fields) > 0 {
			plan.Fields = modelutils.SliceToStringListValue(fields)
		}
	}

	if component.UserConfig["options"] != nil {
		opts, _ := component.UserConfig["options"].(map[string]any)
		if opts["delimiter"] != nil {
			plan.Delimiter = StringValue(string(opts["delimiter"].(string)))
		}
	}
}
