package processors

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	"golang.org/x/exp/slices"
)

var ParseSequentiallyProcessorName = "parse_sequentially"

type ParseSequentiallyProcessorModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Field        String `tfsdk:"field" user_config:"true"`
	TargetField  String `tfsdk:"target_field" user_config:"true"`
	Parsers      List   `tfsdk:"parsers" user_config:"true"`
}

var ParseSequentiallyProcessorResourceSchema = schema.Schema{
	Description: "Parse a field using one of a list of ordered parsers. Parsing ends (short-circuits) on the first successful parse.",
	Attributes:  ExtendBaseAttributes(parse_sequential_schema),
}

func ParseSequentiallyProcessorFromModel(plan *ParseSequentiallyProcessorModel, previousState *ParseSequentiallyProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "parse-sequentially",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"field": plan.Field.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	if !plan.TargetField.IsNull() && !plan.TargetField.IsUnknown() {
		component.UserConfig["target_field"] = plan.TargetField.ValueString()
	}

	parsers, parsers_diag := plan.getParsers()
	dd.Append(parsers_diag...)

	if !parsers_diag.HasError() {
		component.UserConfig["parsers"] = parsers
	}

	return &component, dd
}

func (plan *ParseSequentiallyProcessorModel) getParsers() ([]map[string]any, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	parser_elems := plan.Parsers.Elements()
	var parsers []map[string]any

	for elem_index, elem := range parser_elems {
		parser_item, parser_diag := parserFromModel(elem, elem_index)
		dd.Append(parser_diag...)
		parsers = append(parsers, parser_item)
	}

	return parsers, dd
}

func parserFromModel(elem attr.Value, elem_index int) (map[string]any, diag.Diagnostics) {
	elem_attrs := elem.(Object).Attributes()
	parser := elem_attrs["parser"].(String).ValueString()

	parser_item := map[string]any{}

	if api_parser, ok := VRL_PARSERS[parser]; ok {
		parser_item["parser"] = api_parser
	}

	if !elem_attrs["label"].IsNull() {
		parser_item["label"] = elem_attrs["label"].(String).ValueString()
	}

	options_key := optionsKeyForParser(parser)

	dd := diag.Diagnostics{}
	// options_key for parsers with no options does not exist
	if model_options, ok := elem_attrs[options_key]; ok && !model_options.IsNull() {
		new_options := MapValuesToMapAny(elem_attrs[options_key], &dd)

		if !dd.HasError() {
			parser_item["options"] = new_options
		}
	}

	requires_options := slices.Contains(VRL_PARSERS_WITH_REQUIRED_OPTIONS, parser)
	if _, ok := parser_item["options"]; requires_options && !ok {
		options_key := optionsKeyForParser(parser)
		dd.AddAttributeError(
			path.Root("parsers").AtListIndex(elem_index).AtName(options_key),
			"Incorrect attribute value type",
			fmt.Sprintf("Inappropriate value for attribute \"parsers\": element %v: attribute \"%s\" is required.", elem_index, options_key),
		)
	}

	return parser_item, dd
}

func ParseSequentiallyProcessorToModel(plan *ParseSequentiallyProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Field = StringValue(component.UserConfig["field"].(string))

	if component.UserConfig["target_field"] != nil {
		plan.TargetField = StringValue(component.UserConfig["target_field"].(string))
	}

	api_parsers := parsersFromAPI(component)
	if len(api_parsers) > 0 {
		plan.setParsers(api_parsers)
	}
}

// convert from []any to []map[string]any
func parsersFromAPI(component *Processor) []map[string]any {
	var parsers []map[string]any

	if api_parsers, ok := component.UserConfig["parsers"].([]any); ok {
		for _, item := range api_parsers {
			parsers = append(parsers, item.(map[string]any))
		}
	}
	return parsers
}

func (plan *ParseSequentiallyProcessorModel) setParsers(api_parsers []map[string]any) {
	var parsers []attr.Value

	for _, api_parser_item := range api_parsers {
		parser_item := basetypes.NewObjectValueMust(
			ToAttrTypes(parse_sequential_item_schema),
			parserItemToModel(api_parser_item),
		)

		parsers = append(parsers, parser_item)
	}

	if len(parsers) > 0 {
		plan.Parsers = basetypes.NewListValueMust(plan.Parsers.ElementType(context.Background()), parsers)
	}
}

// convert parser entry in API response to map of Terraform values
func parserItemToModel(api_parser_item map[string]any) map[string]attr.Value {
	parser_item := map[string]attr.Value{}

	parser := FindKey(VRL_PARSERS, api_parser_item["parser"].(string))
	parser_item["parser"] = StringValue(parser)

	if api_parser_item["label"] != nil {
		parser_item["label"] = StringValue(api_parser_item["label"].(string))
	} else {
		parser_item["label"] = basetypes.NewStringNull()
	}

	if api_parser_item["_output_name"] != nil {
		parser_item["output_name"] = StringValue(api_parser_item["_output_name"].(string))
	}

	api_options, ok := api_parser_item["options"].(map[string]any)

	if ok && len(api_options) > 0 {
		options_key := optionsKeyForParser(parser)
		options_attrs, _ := optionsSchemaAttributes(options_key)
		attr_types := ToAttrTypes(options_attrs)
		optional_fields, ok := OPTIONAL_FIELDS_BY_PARSER[parser]
		if !ok {
			optional_fields = []string{}
		}
		values := MapAnyFillMissingValues(attr_types, api_options, optional_fields)
		parser_item[options_key] = basetypes.NewObjectValueMust(attr_types, values)
	}

	// set options for other parsers to null objects to enable conversion of
	// the parser item entry to Terraform model
	for _, key := range MapKeys(base_options_parser_schema) {
		if _, ok := parser_item[key]; ok {
			continue
		}
		schema_attrs, _ := optionsSchemaAttributes(key)
		parser_item[key] = basetypes.NewObjectNull(ToAttrTypes(schema_attrs))
	}

	return parser_item
}
