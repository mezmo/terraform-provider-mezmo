package models

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
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

func PipelineFromModel(plan *PipelineResourceModel) *Pipeline {
	pipeline := Pipeline{
		Title: plan.Title.ValueString(),
	}
	if !plan.Id.IsUnknown() {
		pipeline.Id = plan.Id.ValueString()
	}

	return &pipeline
}

func PipelineToModel(plan *PipelineResourceModel, pipeline *Pipeline) {
	plan.Id = StringValue(pipeline.Id)
	plan.Title = StringValue(pipeline.Title)
	if pipeline.CreatedAt != nil {
		plan.CreatedAt = StringValue(pipeline.CreatedAt.Format(time.RFC3339))
	}
}
