package sources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type HttpSourceModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	Decoding        String `tfsdk:"decoding"`
	CaptureMetadata Bool   `tfsdk:"capture_metadata"`
}

func HttpSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents an HTTP source.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"decoding": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("json"),
				Description: "The decoding method for converting frames into data events.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"bytes", "json", "ndjson"),
				},
			},
		}, true),
	}
}

func HttpSourceFromModel(plan *HttpSourceModel, previousState *HttpSourceModel) *Component {
	component := Component{
		Type:        "http",
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		UserConfig: map[string]any{
			"decoding":         plan.Decoding.ValueString(),
			"capture_metadata": plan.CaptureMetadata.ValueBool(),
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component
}

func HttpSourceToModel(plan *HttpSourceModel, component *Component) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	if component.UserConfig["format"] != nil {
		decoding, _ := component.UserConfig["decoding"].(string)
		plan.Decoding = StringValue(decoding)
		captureMetadata, _ := component.UserConfig["capture_metadata"].(bool)
		plan.CaptureMetadata = BoolValue(captureMetadata)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
}
