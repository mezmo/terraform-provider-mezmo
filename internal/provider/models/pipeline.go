package models

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
)

type PipelineResourceModel struct {
	Id        String `tfsdk:"id"`
	Title     String `tfsdk:"title"`
	CreatedAt String `tfsdk:"created_at"`
	UpdatedAt String `tfsdk:"updated_at"`
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
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func PipelineFromModel(plan *PipelineResourceModel) *Pipeline {
	pipeline := Pipeline{
		Title:  plan.Title.ValueString(),
		Origin: ORIGIN_TERRAFORM,
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
		plan.CreatedAt = StringValue(pipeline.CreatedAt.Format(time.RFC3339Nano))
	}
	if pipeline.UpdatedAt != nil {
		plan.UpdatedAt = StringValue(pipeline.UpdatedAt.Format(time.RFC3339Nano))
	} else {
		plan.UpdatedAt = StringValue(pipeline.CreatedAt.Format(time.RFC3339Nano))
	}
}
