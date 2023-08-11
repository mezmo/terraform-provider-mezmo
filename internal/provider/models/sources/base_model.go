package sources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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

var pushSourceCaptureMetadataBase = SchemaAttributes{
	"capture_metadata": schema.BoolAttribute{
		Computed: true,
		Optional: true,
		Default:  booldefault.StaticBool(false),
		Description: "Enable the inclusion of all http headers and query string parameters " +
			"that were sent from the source",
	},
}

func ExtendBaseAttributes(target SchemaAttributes, is_push_source bool) SchemaAttributes {
	for k, v := range baseSourceSchemaAttributes {
		target[k] = v
	}
	if is_push_source {
		for k, v := range pushSourceCaptureMetadataBase {
			target[k] = v
		}
	}
	return target
}
