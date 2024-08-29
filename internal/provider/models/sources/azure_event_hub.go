package sources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const AZURE_EVENT_HUB_SOURCE_TYPE_NAME = "azure_event_hub"
const AZURE_EVENT_HUB_SOURCE_NODE_NAME = "azure-event-hub"

type AzureEventHubSourceModel struct {
	Id               String `tfsdk:"id"`
	PipelineId       String `tfsdk:"pipeline_id"`
	Title            String `tfsdk:"title"`
	Description      String `tfsdk:"description"`
	GenerationId     Int64  `tfsdk:"generation_id"`
	Decoding         String `tfsdk:"decoding" user_config:"true"`
	ConnectionString String `tfsdk:"connection_string" user_config:"true"`
	Namespace        String `tfsdk:"namespace" user_config:"true"`
	GroupId          String `tfsdk:"group_id" user_config:"true"`
	Topics           List   `tfsdk:"topics" user_config:"true"`
}

var AzureEventHubSourceResourceSchema = schema.Schema{
	Description: "Represents an Azure Event Hub source.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"decoding": schema.StringAttribute{
			Required:    false,
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("bytes"),
			Description: "Configures how events are decoded from raw bytes",
			Validators: []validator.String{
				stringvalidator.OneOf("bytes", "json"),
			},
		},
		"connection_string": schema.StringAttribute{
			Required:    true,
			Description: "The Connection String as it appears in hub consumer SAS Policy",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(512),
			},
		},
		"namespace": schema.StringAttribute{
			Required:    true,
			Description: "The Event Hub Namespace",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(256),
			},
		},
		"group_id": schema.StringAttribute{
			Required:    true,
			Description: "The consumer group name that this consumer belongs to.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(256),
			},
		},
		"topics": schema.ListAttribute{
			Required:    true,
			ElementType: StringType,
			Description: "The list of Azure Event Hub name(s) to read events from.",
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(256),
				),
			},
		},
	}, nil),
}

func AzureEventHubSourceFromModel(plan *AzureEventHubSourceModel, previousState *AzureEventHubSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        AZURE_EVENT_HUB_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"decoding_codec":    plan.Decoding.ValueString(),
				"connection_string": plan.ConnectionString.ValueString(),
				"namespace":         plan.Namespace.ValueString(),
				"group_id":          plan.GroupId.ValueString(),
				"topics":            StringListValueToStringSlice(plan.Topics),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func AzureEventHubSourceToModel(plan *AzureEventHubSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Decoding = StringValue(component.UserConfig["decoding_codec"].(string))
	plan.ConnectionString = StringValue(component.UserConfig["connection_string"].(string))
	plan.Namespace = StringValue(component.UserConfig["namespace"].(string))
	plan.GroupId = StringValue(component.UserConfig["group_id"].(string))
	plan.Topics = SliceToStringListValue(component.UserConfig["topics"].([]any))
}
