package transforms

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type ParseSequentiallyTransformModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Field        String `tfsdk:"field"`
	TargetField  String `tfsdk:"target_field"`
	Parsers      List   `tfsdk:"parsers"`
}

func ParseSequentiallyTransformResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Parse fields using one or more parsers in order, " +
			"sending out the event on the first successful match.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"field": schema.StringAttribute{
				Required:    true,
				Description: "The field whose value should be parsed",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"target_field": schema.StringAttribute{
				Optional: true,
				Description: "The field into which the parsed value should be inserted. Leave blank to " +
					"insert the parsed data into the original field.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
				},
			},
			"parsers": schema.ListNestedAttribute{
				Required: true,
				Description: "The list of parsers to use in order against the input value " +
					"from \"field\", short-circuiting on the first successful match.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label": schema.StringAttribute{
							Optional: true,
							Description: "An arbitrary name you choose to identify the results of this " +
								"processor for use when connecting them to other components.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(20),
							},
						},
						"options": schema.SingleNestedAttribute{
							Attributes:  map[string]schema.Attribute{},
							Optional:    true,
							Description: "Options to pass to the specified parser (if available)",
						},
						"parser": schema.StringAttribute{
							Required:    true,
							Description: "The kind of parser to use against the input value from \"field\".",
							Validators: []validator.String{
								stringvalidator.OneOf(VrlParsers...),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		}),
	}
}

func ParseSequentiallyTransformFromModel(plan *ParseSequentiallyTransformModel, previousState *ParseSequentiallyTransformModel) (*Transform, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Transform{
		BaseNode: BaseNode{
			Type:        "parse-sequentially",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"field": plan.Field.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	p := plan.Parsers.Elements()
	fmt.Printf("----------- uhhhh %+v\n --------------------", p)
	if !plan.TargetField.IsNull() {
		component.UserConfig["target_field"] = plan.TargetField.ValueString()
	}

	return &component, dd
}

func ParseSequentiallyTransformToModel(plan *ParseSequentiallyTransformModel, component *Transform) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
}
