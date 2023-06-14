package models

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client/types"
)

type PipelineResourceModel struct {
	Id        String `tfsdk:"id"`
	Title     String `tfsdk:"title"`
	CreatedAt String `tfsdk:"created_at"`
}

func PipelineResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func PipelineFromModel(model *PipelineResourceModel) *Pipeline {
	pipeline := Pipeline{
		Title: model.Title.ValueString(),
	}
	if !model.Id.IsUnknown() {
		pipeline.Id = model.Id.ValueString()
	}

	return &pipeline
}

func PipelineToModel(model *PipelineResourceModel, pipeline *Pipeline) {
	model.Id = StringValue(pipeline.Id)
	model.Title = StringValue(pipeline.Title)
	if pipeline.CreatedAt != nil {
		model.CreatedAt = StringValue(pipeline.CreatedAt.Format(time.RFC3339))
	}
}
