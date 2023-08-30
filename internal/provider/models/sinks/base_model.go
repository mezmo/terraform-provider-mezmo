package sinks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SchemaAttributes map[string]schema.Attribute

var baseSinkSchemaAttributes = SchemaAttributes{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The uuid of the sink component",
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
		Description: "A user-defined title for the sink component",
	},
	"description": schema.StringAttribute{
		Optional:    true,
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		Description: "A user-defined value describing the sink component",
	},
	"generation_id": schema.Int64Attribute{
		Computed:    true,
		Description: "An internal field used for component versioning",
	},
	"ack_enabled": schema.BoolAttribute{
		Optional:    true,
		Computed:    true,
		Default:     booldefault.StaticBool(true),
		Description: "Acknowledge data from the source when it reaches the sink",
	},
	"inputs": schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true, // The server could set a default, so it's computed
		Description: "The ids of the input components",
	},
}

var addSchemas = map[string]schema.Attribute{
	"batch_timeout_secs": schema.Int64Attribute{
		Computed: true,
		Optional: true,
		Default:  int64default.StaticInt64(300),
		Description: "The maximum amount of time, in seconds, events will be buffered " +
			"before being flushed to the destination",
		Validators: []validator.Int64{int64validator.OneOf(30, 60, 90, 120, 300)},
	},
}

func ExtendBaseAttributes(target SchemaAttributes, addons []string) SchemaAttributes {
	for k, v := range baseSinkSchemaAttributes {
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
