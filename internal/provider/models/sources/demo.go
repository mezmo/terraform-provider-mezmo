package sources

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client/types"
)

type DemoSourceModel struct {
	Id          String `tfsdk:"id"`
	PipelineId  String `tfsdk:"pipeline"`
	Title       String `tfsdk:"title"`
	Description String `tfsdk:"description"`
	Format      String `tfsdk:"format"`
	CreatedAt   String `tfsdk:"created_at"`
}

func DemoSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pipeline": schema.StringAttribute{
				Required: true,
			},
			"title": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
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
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func DemoSourceFromModel(model *DemoSourceModel) *Component {
	component := Component{
		Title:       model.Title.ValueString(),
		Description: model.Format.ValueString(),
		UserOptions: map[string]any{"format": model.Format.ValueString()},
	}
	if !model.Id.IsUnknown() {
		component.Id = model.Id.ValueString()
	}

	return &component
}

func DemoSourceToModel(model *DemoSourceModel, component *Component) {
	model.Id = StringValue(component.Id)
	model.Title = StringValue(component.Title)
	model.Description = StringValue(component.Description)
	if component.UserOptions["string"] != nil {
		model.Format = StringValue(component.UserOptions["string"].(string))
	}
	if component.CreatedAt != nil {
		model.CreatedAt = StringValue(component.CreatedAt.Format(time.RFC3339))
	}
}
