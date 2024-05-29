package models

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
)

var PUSH_SOURCE_TYPES = []string{
	"http",
	"splunk-hec",
	"kinesis-firehose",
	"fluent",
	"logstash",
	"mezmo-agent",
	"mezmo-datadog-source",
	"webhook",
	"prometheus-remote-write",
	"open-telemetry-metrics",
	"open-telemetry-logs",
	"open-telemetry-traces",
}

type SharedSourceResourceModel struct {
	Id          StringValue `tfsdk:"id"`
	Title       StringValue `tfsdk:"title" user_config:"true"`
	Description StringValue `tfsdk:"description" user_config:"true"`
	Type        StringValue `tfsdk:"type" user_config:"true"`
}

func SharedSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the shared source.",
				Computed:    true,
			},
			"title": schema.StringAttribute{
				Description: "A descriptive name for the shared source.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
				},
			},
			"description": schema.StringAttribute{
				Description: "Details describing the shared source.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of source that should be shared. This is typically the " +
					"name of the source componenet, e.g. `kinesis-firehose`. The source must be a " +
					"\"push source\", thus pull sources cannot be shared and will result in an error.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(PUSH_SOURCE_TYPES...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// From terraform schema/model to a struct for sending to the API
func SharedSourceFromModel(plan *SharedSourceResourceModel) *SharedSource {
	source := SharedSource{
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Type:        plan.Type.ValueString(),
	}
	if !plan.Id.IsUnknown() {
		source.Id = plan.Id.ValueString()
	}
	return &source
}

// From an API response to a terraform model
func SharedSourceToModel(plan *SharedSourceResourceModel, source *SharedSource) {
	plan.Id = NewStringValue(source.Id)
	plan.Title = NewStringValue(source.Title)
	plan.Description = NewStringValue(source.Description)
	plan.Type = NewStringValue(source.Type)
}
