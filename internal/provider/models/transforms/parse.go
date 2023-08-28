package transforms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
	"golang.org/x/exp/slices"
)

type ParseTransformModel struct {
	Id               String `tfsdk:"id"`
	PipelineId       String `tfsdk:"pipeline_id"`
	Title            String `tfsdk:"title"`
	Description      String `tfsdk:"description"`
	Inputs           List   `tfsdk:"inputs"`
	GenerationId     Int64  `tfsdk:"generation_id"`
	Field            String `tfsdk:"field"`
	TargetField      String `tfsdk:"target_field"`
	Parser           String `tfsdk:"parser"`
	ApacheOptions    Object `tfsdk:"apache_log_options"`
	CefOptions       Object `tfsdk:"cef_log_options"`
	CsvOptions       Object `tfsdk:"csv_row_options"`
	GrokOptions      Object `tfsdk:"grok_parser_options"`
	KeyValueOptions  Object `tfsdk:"key_value_log_options"`
	NginxOptions     Object `tfsdk:"nginx_log_options"`
	RegexOptions     Object `tfsdk:"regex_parser_options"`
	TimestampOptions Object `tfsdk:"timestamp_parser_options"`
}

func ParseTransformResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Parse a specified field using the chosen parser",
		Attributes:  ExtendBaseAttributes(parse_schema),
	}
}

func ParseTransformFromModel(plan *ParseTransformModel, previousState *ParseTransformModel) (*Transform, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	parser := plan.Parser.ValueString()
	component := Transform{
		BaseNode: BaseNode{
			Type:        "parse",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"field": plan.Field.ValueString(),
			},
		},
	}

	if api_parser, ok := VRL_PARSERS[parser]; ok {
		component.UserConfig["parser"] = api_parser
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	if !plan.TargetField.IsNull() {
		component.UserConfig["target_field"] = plan.TargetField.ValueString()
	}

	model_options := plan.optionsFromModel()

	if !model_options.IsNull() {
		options, has_error := MapValuesToMapStrings(model_options, dd)
		if !has_error {
			component.UserConfig["options"] = options
		}
	} else if slices.Contains(VRL_PARSES_WITH_REQUIRED_OPTIONS, parser) {
		options_key := fmt.Sprintf("%s_options", parser)
		dd.AddAttributeError(
			path.Root(options_key),
			fmt.Sprintf("Attribute %s is required.", options_key),
			fmt.Sprintf("Attribute %s is required for %s.", options_key, parser),
		)
	}

	return &component, dd
}

func ParseTransformToModel(plan *ParseTransformModel, component *Transform) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	parser := FindKey(VRL_PARSERS, component.UserConfig["parser"].(string))
	plan.Field = StringValue(component.UserConfig["field"].(string))
	plan.Parser = StringValue(parser)

	if component.UserConfig["target_field"] != nil {
		plan.TargetField = StringValue(component.UserConfig["target_field"].(string))
	}

	options, ok := component.UserConfig["options"].(map[string]any)
	if ok {
		plan.setOptions(options)
	}
}

func (plan *ParseTransformModel) optionsFromModel() basetypes.ObjectValue {
	if !plan.ApacheOptions.IsNull() {
		return plan.ApacheOptions
	}
	if !plan.CefOptions.IsNull() {
		return plan.CefOptions
	}
	if !plan.CsvOptions.IsNull() {
		return plan.CsvOptions
	}
	if !plan.GrokOptions.IsNull() {
		return plan.GrokOptions
	}
	if !plan.KeyValueOptions.IsNull() {
		return plan.KeyValueOptions
	}
	if !plan.NginxOptions.IsNull() {
		return plan.NginxOptions
	}
	if !plan.RegexOptions.IsNull() {
		return plan.RegexOptions
	}
	if !plan.TimestampOptions.IsNull() {
		return plan.TimestampOptions
	}

	return basetypes.NewObjectNull(plan.ApacheOptions.AttributeTypes(context.Background()))
}

func (plan *ParseTransformModel) setOptions(options map[string]any) {
	if len(options) == 0 {
		return
	}

	parser := plan.Parser.ValueString()
	optional_fields, ok := OPTIONAL_FIELDS_BY_PARSER[parser]

	if !ok {
		optional_fields = []string{}
	}

	switch parser {
	case VRL_PARSER_APACHE:
		plan.ApacheOptions = optionsToModel(plan.ApacheOptions, options, optional_fields)
	case VRL_PARSER_CEF:
		plan.CefOptions = optionsToModel(plan.CefOptions, options, optional_fields)
	case VRL_PARSER_CSV:
		plan.CsvOptions = optionsToModel(plan.CsvOptions, options, optional_fields)
	case VRL_PARSER_GROK:
		plan.GrokOptions = optionsToModel(plan.GrokOptions, options, optional_fields)
	case VRL_PARSER_KEY_VALUE:
		plan.KeyValueOptions = optionsToModel(plan.KeyValueOptions, options, optional_fields)
	case VRL_PARSER_NGINX:
		plan.NginxOptions = optionsToModel(plan.NginxOptions, options, optional_fields)
	case VRL_PARSER_REGEX:
		plan.RegexOptions = optionsToModel(plan.RegexOptions, options, optional_fields)
	case VRL_PARSER_TIMESTAMP:
		plan.TimestampOptions = optionsToModel(plan.TimestampOptions, options, optional_fields)
	}
}

func optionsToModel(obj basetypes.ObjectValue, options map[string]any, optional_fields []string) basetypes.ObjectValue {
	attr_types := obj.AttributeTypes(context.Background())

	if len(options) > 0 {
		values := MapAnyToMapValues(attr_types, options, optional_fields)
		new_options := basetypes.NewObjectValueMust(attr_types, values)

		return new_options
	}

	return basetypes.NewObjectNull(attr_types)
}
