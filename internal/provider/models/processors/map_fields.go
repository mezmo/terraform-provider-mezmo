package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const MAP_FIELDS_PROCESSOR_NODE_NAME = "map-fields"
const MAP_FIELDS_PROCESSOR_TYPE_NAME = "map_fields"

type MapFieldsProcessorModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	Inputs       ListValue   `tfsdk:"inputs"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
	Mappings     ListValue   `tfsdk:"mappings" user_config:"true"`
}

var mapFieldsAttrTypes = map[string]attr.Type{
	"source_field":     StringType{},
	"target_field":     StringType{},
	"drop_source":      BoolType{},
	"overwrite_target": BoolType{},
}

var MapFieldsProcessorResourceSchema = schema.Schema{
	Description: "Maps data from one field to another, either by moving or copying",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"mappings": schema.ListNestedAttribute{
			Required:    true,
			Description: "A list of field mappings. Mappings are applied in the order they are defined",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"source_field": schema.StringAttribute{
						Required:    true,
						Description: "The field to copy data from",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"target_field": schema.StringAttribute{
						Required:    true,
						Description: "The field to copy data into",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"drop_source": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
						Description: "When enabled, the source field is dropped after the data is copied " +
							"to the target field. Otherwise, it is preserved.",
					},
					"overwrite_target": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
						Description: "When enabled, any existing data in the target field is overwritten. " +
							"Otherwise, the target field will be preserved and this mapping will " +
							"have no effect.",
					},
				},
			},
		},
	}),
}

func MapFieldsProcessorFromModel(plan *MapFieldsProcessorModel, previousState *MapFieldsProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := &Processor{
		BaseNode: BaseNode{
			Type:        MAP_FIELDS_PROCESSOR_NODE_NAME,
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

	if !plan.Mappings.IsNull() {
		mappings := make([]map[string]any, 0)
		for _, v := range plan.Mappings.Elements() {
			obj := MapValuesToMapAny(v, &dd)
			mappings = append(mappings, obj)
		}
		component.UserConfig["mappings"] = mappings
	}

	return component, dd
}

func MapFieldsProcessorToModel(plan *MapFieldsProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	if component.UserConfig["mappings"] != nil {
		mappings := make([]attr.Value, 0)
		for _, v := range component.UserConfig["mappings"].([]any) {
			values := v.(map[string]any)
			attrValues := map[string]attr.Value{
				"source_field":     NewStringValue(values["source_field"].(string)),
				"target_field":     NewStringValue(values["target_field"].(string)),
				"drop_source":      NewBoolValue(values["drop_source"].(bool)),
				"overwrite_target": NewBoolValue(values["overwrite_target"].(bool)),
			}
			mappings = append(mappings, NewObjectValueMust(mapFieldsAttrTypes, attrValues))

			plan.Mappings = NewListValueMust(ObjectType{AttrTypes: mapFieldsAttrTypes}, mappings)
		}
	}
}
