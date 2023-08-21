package sources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type DemoSourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Format       String `tfsdk:"format"`
	GenerationId Int64  `tfsdk:"generation_id"`
}

func DemoSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents a demo logs source.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The format of the events",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"env_sensor", "financial", "nginx", "json", "apache_common",
						"apache_error", "bsd_syslog", "syslog", "http_metrics", "generic_metrics"),
				},
			},
		}, nil),
	}
}

func DemoSourceFromModel(plan *DemoSourceModel, previousState *DemoSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Source{
		BaseNode: BaseNode{
			Type:        "demo-logs",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{"format": plan.Format.ValueString()},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func DemoSourceToModel(plan *DemoSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	if component.UserConfig["format"] != nil {
		format, _ := component.UserConfig["format"].(string)
		plan.Format = StringValue(format)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
}
