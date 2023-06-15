package sources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type DemoSourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Format       String `tfsdk:"format"`
	GenerationId Int64  `tfsdk:"generation_id"`
}

func DemoSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pipeline": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"title": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The format of the events",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"env_sensor", "financial", "nginx", "json", "apache_common",
						"apache_error", "bsd_syslog", "syslog", "http_metrics", "generic_metrics"),
				},
			},
			"generation_id": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func DemoSourceFromModel(model *DemoSourceModel, previousState *DemoSourceModel) *Component {
	component := Component{
		Type:        "demo-logs",
		Title:       model.Title.ValueString(),
		Description: model.Description.ValueString(),
		UserConfig:  map[string]any{"format": model.Format.ValueString()},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component
}

func DemoSourceToModel(model *DemoSourceModel, component *Component) {
	model.Id = StringValue(component.Id)
	if component.Title != "" {
		model.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		model.Description = StringValue(component.Description)
	}
	if component.UserConfig["format"] != nil {
		format, _ := component.UserConfig["format"].(string)
		model.Format = StringValue(format)
	}
	model.GenerationId = Int64Value(component.GenerationId)
}
