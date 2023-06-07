package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	. "github.com/hashicorp/terraform-plugin-framework/types"
)

type pipelineResourceModel struct {
	Id        String `tfsdk:"id"`
	Title     String `tfsdk:"title"`
	UpdatedAt String `tfsdk:"updated_at"`
}

func pipelineResourceSchema() schema.Schema {
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
