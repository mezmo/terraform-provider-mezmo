package models

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
)

type PublishPipelineResourceModel struct {
	PipelineId StringValue `tfsdk:"pipeline_id"`
}

func PublishPipelineResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource will monitor a pipeline for changes, and publish it when necessary.\n" +
			"\n## Configuration\n" +
			"To make sure a pipeline and its components exist before publishing, the configuration of this resource " +
			"requires the use of child modules and `depends_on`. The pipeline's configuration should exist in a " +
			"child module with an `output` of the pipeline's `id` field. This resource will then reference " +
			"this field as `pipeline_id`, and be able to publish as needed when the pipeline changes.\n" +
			"\nThe `output` can be done however the user chooses, as long as the pipeline's `id` is accessible in the root module. " +
			"In other words, `output` can be an object of the entire pipeline, or just the `id`.\n" +
			"\n## Example Child Module\n" +
			"```terraform\n" +
			`
# This would exist in an arbitrary module directory. For example, "./modules/main.tf"
terraform {
	required_providers {
		mezmo = {
			source = "registry.terraform.io/mezmo/mezmo"
		}
	}
}
output "my_pipeline" {
	value = mezmo_pipeline.my_pipeline
}
resource "mezmo_pipeline" "my_pipeline" {
	title = "A pipeline to publish"
}
# ... other sources, processors, destinations, etc.
			` +
			"\n```\n",
		Attributes: map[string]schema.Attribute{
			"pipeline_id": schema.StringAttribute{
				Description: "The id of the pipeline to monitor for publishing. " +
					"Any changes to its components will trigger a publish. " +
					"This pipeline must be configured in a child module with an `output`.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// From terraform schema/model to a struct for sending to the API
func PublishPipelineFromModel(plan *PublishPipelineResourceModel) *PublishPipeline {
	publishPipeline := PublishPipeline{
		PipelineId: plan.PipelineId.ValueString(),
	}
	return &publishPipeline
}

// From an API response to a terraform model
func PublishPipelineToModel(plan *PublishPipelineResourceModel, publishPipeline *PublishPipeline) {
	plan.PipelineId = NewStringValue(publishPipeline.PipelineId)
}
