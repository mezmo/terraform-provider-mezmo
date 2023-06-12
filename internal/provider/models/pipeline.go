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
	UpdatedAt String `tfsdk:"updated_at"`
}

func PipelineResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Optional: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func PipelineFromModel(model *PipelineResourceModel) *Pipeline {
	pipeline := Pipeline {
		Title:     model.Title.ValueString(),
	}
	if !model.Id.IsUnknown() {
		pipeline.Id = model.Id.ValueString()
	}
	if !model.UpdatedAt.IsUnknown() {
		updatedAt, _ := time.Parse(time.RFC3339, model.UpdatedAt.ValueString())
		pipeline.UpdatedAt = updatedAt
	}

	return &pipeline
}

func PipelineToModel(model *PipelineResourceModel, pipeline *Pipeline) {
	model.Id = StringValue(pipeline.Id)
	model.Title = StringValue(pipeline.Title)
	model.UpdatedAt = StringValue(pipeline.UpdatedAt.Format(time.RFC3339))
}
