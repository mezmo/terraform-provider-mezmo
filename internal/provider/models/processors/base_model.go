package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SchemaAttributes map[string]schema.Attribute

var baseProcessorSchemaAttributes = SchemaAttributes{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The uuid of the processor",
	},
	"pipeline_id": schema.StringAttribute{
		Required:    true,
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		Description: "The uuid of the pipeline",
	},
	"title": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
			stringvalidator.LengthAtMost(256),
		},
		Description: "A user-defined title for the processor",
	},
	"description": schema.StringAttribute{
		Optional:    true,
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		Description: "A user-defined value describing the processor",
	},
	"generation_id": schema.Int64Attribute{
		Computed:    true,
		Description: "An internal field used for component versioning",
	},
	"inputs": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true, // The server could set a default, so it's computed
		Description: "The ids of the input components",
	},
}

func ExtendBaseAttributes(target SchemaAttributes) SchemaAttributes {
	return ExtendSchemaAttributes(baseProcessorSchemaAttributes, target)
}

func ExtendSchemaAttributes(fromAttributes SchemaAttributes, toAttributes SchemaAttributes) SchemaAttributes {
	for k, v := range fromAttributes {
		toAttributes[k] = v
	}
	return toAttributes
}
