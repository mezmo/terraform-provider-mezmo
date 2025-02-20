package sources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type SchemaAttributes map[string]schema.Attribute

var baseSourceSchemaAttributes = SchemaAttributes{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The uuid of the source component",
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
		Description: "A user-defined title for the source component",
	},
	"description": schema.StringAttribute{
		Optional:    true,
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		Description: "A user-defined value describing the source component",
	},
	"generation_id": schema.Int64Attribute{
		Computed:    true,
		Description: "An internal field used for component versioning",
	},
}

var addSchemas = map[string]schema.Attribute{
	"shared_source_id": schema.StringAttribute{
		Computed: true,
		Optional: true,
		Description: "The uuid of a pipeline source or shared source to be used as the input for this " +
			"component. This can only be provided on resource creation (not update).",
	},
}

func ExtendBaseAttributes(target SchemaAttributes, addons []string) SchemaAttributes {
	for k, v := range baseSourceSchemaAttributes {
		target[k] = v
	}
	for _, name := range addons {
		schema, ok := addSchemas[name]
		if !ok {
			panic(fmt.Errorf("Addon attribute %s not found. Developer error.", name))
		}
		target[name] = schema
	}
	return target
}
