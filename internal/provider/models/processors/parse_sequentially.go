package processors

import (
	"context"
	"fmt"
	"strings"

	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const PARSE_SEQUENTIALLY_PROCESSOR_TYPE_NAME = "parse_sequentially"
const PARSE_SEQUENTIALLY_PROCESSOR_NODE_NAME = "parse-sequentially"

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
	// unmatched exists in outputs from API which is not exposed to the
	// user. Mapped to user_config to make TF happy
	Unmatched String `tfsdk:"unmatched" user_config:"true"`
}

var ParseSequentiallyProcessorResourceSchema = schema.Schema{
	Description: "Parse a field using one of a list of ordered parsers. Parsing ends (short-circuits) on the first successful parse.",
	Attributes:  ExtendBaseAttributes(parse_sequential_schema),
}

func ParseSequentiallyProcessorFromModel(plan *ParseSequentiallyProcessorModel, previousState *ParseSequentiallyProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        PARSE_SEQUENTIALLY_PROCESSOR_NODE_NAME,
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
		parserItem, parser_diag := parserFromModel(elem, elem_index)
		dd.Append(parser_diag...)
		parsers = append(parsers, parserItem)
	}

	return parsers, dd
}

func parserFromModel(elem attr.Value, elem_index int) (map[string]any, diag.Diagnostics) {
	elem_attrs := elem.(Object).Attributes()
	parser := elem_attrs["parser"].(String).ValueString()

	parserItem := map[string]any{}

	if api_parser, ok := VRL_PARSERS[parser]; ok {
		parserItem["parser"] = api_parser
	}

	if !elem_attrs["label"].IsNull() {
		parserItem["label"] = elem_attrs["label"].(String).ValueString()
	}

	optionsKey := optionsKeyForParser(parser)

	dd := diag.Diagnostics{}
	// optionsKey for parsers with no options does not exist
	if model_options, ok := elem_attrs[optionsKey]; ok && !model_options.IsNull() {
		new_options := MapValuesToMapAny(elem_attrs[optionsKey], &dd)

		if !dd.HasError() {
			parserItem["options"] = new_options
		}
	}

	requires_options := slices.Contains(VRL_PARSERS_WITH_REQUIRED_OPTIONS, parser)
	if _, ok := parserItem["options"]; requires_options && !ok {
		optionsKey := optionsKeyForParser(parser)
		dd.AddAttributeError(
			path.Root("parsers").AtListIndex(elem_index).AtName(optionsKey),
			"Incorrect attribute value type",
			fmt.Sprintf("Inappropriate value for attribute \"parsers\": element %v: attribute \"%s\" is required.", elem_index, optionsKey),
		)
	}

	return parserItem, dd
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

	outputNames := mappedOutputs(component)
	api_parsers := parsersFromAPI(component)
	if len(api_parsers) > 0 {
		plan.setParsers(api_parsers, outputNames)
	}

	// Set the unmatched output
	if unmatched, ok := outputNames["_unmatched"]; ok {
		plan.Unmatched = unmatched
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

func (plan *ParseSequentiallyProcessorModel) setParsers(api_parsers []map[string]any, outputNames map[string]basetypes.StringValue) {
	var parsers []attr.Value

	for _, apiParserItem := range api_parsers {
		parserItem := basetypes.NewObjectValueMust(
			ToAttrTypes(parse_sequential_item_schema),
			parserItemToModel(apiParserItem, outputNames),
		)

		parsers = append(parsers, parserItem)
	}

	if len(parsers) == 0 {
		return
	}
	elemType := plan.Parsers.ElementType(context.Background())
	// used by ConvertToTerraformModel method
	// when hydrating without a terraform state, the list element type is nil
	if elemType == nil {
		listType := ParseSequentiallyProcessorResourceSchema.Attributes["parsers"].GetType()
		elemType = listType.(basetypes.ListType).ElementType()
	}
	plan.Parsers = basetypes.NewListValueMust(elemType, parsers)
}

// convert parser entry in API response to map of Terraform values
func parserItemToModel(apiParserItem map[string]any, outputNames map[string]basetypes.StringValue) map[string]attr.Value {
	parserItem := map[string]attr.Value{}

	parser := FindKey(VRL_PARSERS, apiParserItem["parser"].(string))
	parserItem["parser"] = StringValue(parser)

	if apiParserItem["label"] != nil {
		parserItem["label"] = StringValue(apiParserItem["label"].(string))
	} else {
		parserItem["label"] = basetypes.NewStringNull()
	}

	if apiParserItem["_output_name"] != nil {
		apiOutputName := strings.ToLower(apiParserItem["_output_name"].(string))
		// use the fully qualified output name
		if parserOutputName, ok := outputNames[apiOutputName]; ok {
			parserItem["output_name"] = parserOutputName
		}
	}

	api_options, ok := apiParserItem["options"].(map[string]any)

	if ok && len(api_options) > 0 {
		optionsKey := optionsKeyForParser(parser)
		optionsAttrs, _ := optionsSchemaAttributes(optionsKey)
		attrTypes := ToAttrTypes(optionsAttrs)
		optionalFields, ok := OPTIONAL_FIELDS_BY_PARSER[parser]
		if !ok {
			optionalFields = []string{}
		}
		attrTypeKeys := append(MapKeys(attrTypes), optionalFields...)
		values := MapAnyFillMissingValues(attrTypes, StripUnknownOptions(attrTypeKeys, api_options), optionalFields)
		parserItem[optionsKey] = basetypes.NewObjectValueMust(attrTypes, values)
	}

	// set options for other parsers to null objects to enable conversion of
	// the parser item entry to Terraform model
	for _, key := range MapKeys(base_options_parser_schema) {
		if _, ok := parserItem[key]; ok {
			continue
		}
		schema_attrs, _ := optionsSchemaAttributes(key)
		parserItem[key] = basetypes.NewObjectNull(ToAttrTypes(schema_attrs))
	}

	return parserItem
}

func mappedOutputs(component *Processor) map[string]basetypes.StringValue {
	res := make(map[string]basetypes.StringValue, len(component.Outputs))
	for _, item := range component.Outputs {
		idParts := strings.Split(item.Id, ".")
		if len(idParts) == 2 {
			res[strings.ToLower(idParts[1])] = StringValue(item.Id)
		}
	}
	return res
}
