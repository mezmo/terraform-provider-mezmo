package processors

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type EventToMetricProcessorModel struct {
	Id             StringValue  `tfsdk:"id"`
	PipelineId     StringValue  `tfsdk:"pipeline_id"`
	Title          StringValue  `tfsdk:"title"`
	Description    StringValue  `tfsdk:"description"`
	Inputs         ListValue    `tfsdk:"inputs"`
	GenerationId   Int64Value   `tfsdk:"generation_id"`
	MetricName     StringValue  `tfsdk:"metric_name" user_config:"true"`
	MetricKind     StringValue  `tfsdk:"metric_kind" user_config:"true"`
	MetricType     StringValue  `tfsdk:"metric_type" user_config:"true"`
	ValueField     StringValue  `tfsdk:"value_field"`  // Differs from user_config. Manually construct to make types easier.
	ValueNumber    Float64Value `tfsdk:"value_number"` // Ditto. This keeps data types for numbers and strings simpler.
	NamespaceField StringValue  `tfsdk:"namespace_field"`
	NamespaceValue StringValue  `tfsdk:"namespace_value"`
	Tags           ListValue    `tfsdk:"tags" user_config:"true"`
}

var METRIC_NAME_REGEX = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_:]*$")
var METRIC_TAG_NAME_REGEX = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

var metricTagAttrTypes = map[string]attr.Type{
	"name":       StringType{},
	"value_type": StringType{},
	"value":      StringType{},
}

var EventToMetricProcessorResourceSchema = schema.Schema{
	Description: "Allows conversion between arbitrary events and a Metric",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"metric_name": schema.StringAttribute{
			Required:    true,
			Description: "The machine name of the metric to emit",
			Validators: []validator.String{
				stringvalidator.RegexMatches(METRIC_NAME_REGEX, "has invalid characters; See documention"),
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(128),
			},
		},
		"metric_kind": schema.StringAttribute{
			Required: true,
			Description: "The kind of metric to emit, Absolute or Incremental. Absolute metrics represent " +
				"a complete value, and will generally replace an existing value for the metric in " +
				"the target destination. Incremental metrics represent an additive value which " +
				"is aggregated in the target destination to produce a new value.",
			Validators: []validator.String{
				stringvalidator.OneOf(MetricKind...),
			},
		},
		"metric_type": schema.StringAttribute{
			Required:    true,
			Description: "The type of metric to emit. For example, counter, sum, gauge.",
			Validators: []validator.String{
				stringvalidator.OneOf(MetricType...),
			},
		},
		"value_field": schema.StringAttribute{
			Optional:    true,
			Description: "The value of the metric should come from this event field path.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("value_number")),
			},
		},
		"value_number": schema.Float64Attribute{
			Optional:    true,
			Description: "Use this specified numeric value.",
		},
		"namespace_field": schema.StringAttribute{
			Optional:    true,
			Description: "The value of the namespace should come from this event field path.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("namespace_value")),
			},
		},
		"namespace_value": schema.StringAttribute{
			Optional:    true,
			Description: "The namespace value should be this specified string.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"tags": schema.ListNestedAttribute{
			Optional:    true,
			Description: "A set of tags (also called labels) to apply to the metric event.",
			Validators: []validator.List{
				listvalidator.SizeAtMost(10),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The tag name",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.LengthAtMost(128),
							stringvalidator.RegexMatches(METRIC_TAG_NAME_REGEX, "has invalid characters; See documention"),
						},
					},
					"value_type": schema.StringAttribute{
						Required:    true,
						Description: "Specifies if the value comes from an event field, or a new value input.",
						Validators: []validator.String{
							stringvalidator.OneOf("field", "value"),
						},
					},
					"value": schema.StringAttribute{
						Required: true,
						Description: "For value types, this is the value of the tag. If using a field type, the " +
							"value comes from this field path. Note that fields with highly-variable values will result " +
							"in high-cardinality metrics, which may impact storage or cost in downstream destinations.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		},
	}),
}

func EventToMetricProcessorFromModel(plan *EventToMetricProcessorModel, previousState *EventToMetricProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := &Processor{
		BaseNode: BaseNode{
			Type:        "event-to-metric",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"metric_name": plan.MetricName.ValueString(),
				"metric_kind": plan.MetricKind.ValueString(),
				"metric_type": plan.MetricType.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	component.Inputs = StringListValueToStringSlice(plan.Inputs)
	user_config := component.UserConfig

	if plan.ValueField.IsNull() && plan.ValueNumber.IsNull() {
		dd.AddError("Invalid configuration", "Either `value_field` or `value_number` is required")
		return nil, dd
	}
	// Create the proper API request body manually. This reduces type complexity in the provider.
	if !plan.ValueField.IsNull() {
		user_config["value"] = map[string]any{
			"value_type": "field",
			"value":      plan.ValueField.ValueString(),
		}
	}
	if !plan.ValueNumber.IsNull() {
		user_config["value"] = map[string]any{
			"value_type": "value",
			"value":      plan.ValueNumber.ValueFloat64(),
		}
	}

	namespace_map := make(map[string]any)

	if !plan.NamespaceValue.IsNull() {
		namespace_map["value_type"] = "value"
		namespace_map["value"] = plan.NamespaceValue.ValueString()
	} else if !plan.NamespaceField.IsNull() {
		namespace_map["value_type"] = "field"
		namespace_map["value"] = plan.NamespaceField.ValueString()
	} else {
		// namespace can be nullified
		namespace_map["value_type"] = "none"
	}
	user_config["namespace"] = namespace_map

	if !plan.Tags.IsNull() {
		tags := make([]map[string]any, 0)
		for _, v := range plan.Tags.Elements() {
			obj := MapValuesToMapAny(v, &dd)
			tags = append(tags, obj)
		}
		user_config["tags"] = tags
	}

	return component, dd
}

func EventToMetricProcessorToModel(plan *EventToMetricProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	plan.MetricName = NewStringValue(component.UserConfig["metric_name"].(string))
	plan.MetricKind = NewStringValue(component.UserConfig["metric_kind"].(string))
	plan.MetricType = NewStringValue(component.UserConfig["metric_type"].(string))

	value, _ := component.UserConfig["value"].(map[string]any)
	if value["value_type"].(string) == "field" {
		plan.ValueField = NewStringValue(value["value"].(string))
	} else {
		// it's a number
		plan.ValueNumber = NewFloat64Value(value["value"].(float64))
	}

	// Initialize namespace fields with the value_type == "none" case
	plan.NamespaceField = NewStringNull()
	plan.NamespaceValue = NewStringNull()

	namespace, _ := component.UserConfig["namespace"].(map[string]any)
	value_type := namespace["value_type"].(string)

	if value_type == "field" {
		plan.NamespaceField = NewStringValue(namespace["value"].(string))
	} else if value_type == "value" {
		plan.NamespaceValue = NewStringValue(namespace["value"].(string))
	}

	if component.UserConfig["tags"] != nil {
		tags := make([]attr.Value, 0)
		for _, v := range component.UserConfig["tags"].([]any) {
			values := v.(map[string]any)
			attrValues := map[string]attr.Value{
				"name":       NewStringValue(values["name"].(string)),
				"value_type": NewStringValue(values["value_type"].(string)),
				"value":      NewStringValue(values["value"].(string)),
			}
			tags = append(tags, NewObjectValueMust(metricTagAttrTypes, attrValues))
			plan.Tags = NewListValueMust(ObjectType{AttrTypes: metricTagAttrTypes}, tags)
		}
	}
}
