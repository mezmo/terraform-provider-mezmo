package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const GCP_CLOUD_PUBSUB_DESTINATION_TYPE_NAME = "gcp_cloud_pubsub"
const GCP_CLOUD_PUBSUB_DESTINATION_NODE_NAME = "gcp-cloud-pubsub"

type GcpCloudPubSubDestinationModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	Inputs          List   `tfsdk:"inputs"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	Encoding        String `tfsdk:"encoding" user_config:"true"`
	ProjectId       String `tfsdk:"project_id" user_config:"true"`
	Topic           String `tfsdk:"topic" user_config:"true"`
	CredentialsJSON String `tfsdk:"credentials_json" user_config:"true"`
	AckEnabled      Bool   `tfsdk:"ack_enabled" user_config:"true"`
}

var GcpCloudPubSubResourceSchema = schema.Schema{
	Description: "Publish events to GCP Cloud PubSub",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"encoding": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Dictates how the data will be serialized before storing.",
			Default:     stringdefault.StaticString("text"),
			Validators: []validator.String{
				stringvalidator.OneOf("json", "text"),
			},
		},
		"project_id": schema.StringAttribute{
			Required:    true,
			Description: "The Project ID as defined in Google Cloud.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"topic": schema.StringAttribute{
			Required:    true,
			Description: "The name of the topic in which to publish messages.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"credentials_json": schema.StringAttribute{
			Required:    true,
			Description: "JSON Credentials",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}, nil),
}

func GcpCloudPubSubDestinationFromModel(plan *GcpCloudPubSubDestinationModel, previousState *GcpCloudPubSubDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        GCP_CLOUD_PUBSUB_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":      plan.AckEnabled.ValueBool(),
				"encoding":         plan.Encoding.ValueString(),
				"project_id":       plan.ProjectId.ValueString(),
				"topic":            plan.Topic.ValueString(),
				"credentials_json": plan.CredentialsJSON.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func GcpCloudPubSubDestinationToModel(plan *GcpCloudPubSubDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.ProjectId = StringValue(component.UserConfig["project_id"].(string))
	plan.Topic = StringValue(component.UserConfig["topic"].(string))
	plan.CredentialsJSON = StringValue(component.UserConfig["credentials_json"].(string))

}
