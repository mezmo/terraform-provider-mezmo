package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const GCP_CLOUD_MONITORING_DESTINATION_TYPE_NAME = "gcp_cloud_monitoring"
const GCP_CLOUD_MONITORING_DESTINATION_NODE_NAME = "gcp-cloud-monitoring"

type GcpCloudMonitoringDestinationModel struct {
	Id              StringValue `tfsdk:"id"`
	PipelineId      StringValue `tfsdk:"pipeline_id"`
	Title           StringValue `tfsdk:"title"`
	Description     StringValue `tfsdk:"description"`
	Inputs          ListValue   `tfsdk:"inputs"`
	GenerationId    Int64Value  `tfsdk:"generation_id"`
	AckEnabled      BoolValue   `tfsdk:"ack_enabled" user_config:"true"`
	CredentialsJSON StringValue `tfsdk:"credentials_json" user_config:"true"`
	ProjectId       StringValue `tfsdk:"project_id" user_config:"true"`
	ResourceType    StringValue `tfsdk:"resource_type" user_config:"true"`
	ResourceLabels  MapValue    `tfsdk:"resource_labels" user_config:"true"`
}

var GcpCloudMonitoringResourceSchema = schema.Schema{
	Description: "Publish metrics events to GCP Cloud Monitoring",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"credentials_json": schema.StringAttribute{
			Required:    true,
			Description: "JSON Credentials",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"project_id": schema.StringAttribute{
			Required:    true,
			Description: "The Project ID as defined in Google Cloud.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"resource_type": schema.StringAttribute{
			Required:    true,
			Description: "The monitored-resource type as defined in Monitoring.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"resource_labels": schema.MapAttribute{
			Optional:    true,
			ElementType: StringType{},
			Description: "Key/Value pair used to describe the resource",
			Validators: []validator.Map{
				mapvalidator.All(
					mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				),
			},
		},
	}, nil),
}

func GcpCloudMonitoringDestinationFromModel(plan *GcpCloudMonitoringDestinationModel, previousState *GcpCloudMonitoringDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        GCP_CLOUD_MONITORING_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":      plan.AckEnabled.ValueBool(),
				"credentials_json": plan.CredentialsJSON.ValueString(),
				"project_id":       plan.ProjectId.ValueString(),
				"resource_type":    plan.ResourceType.ValueString(),
			},
		},
	}

	if !plan.ResourceLabels.IsNull() {
		labelsMap := MapValuesToMapAny(plan.ResourceLabels, &dd)
		if !dd.HasError() {
			labelsArray := make([]map[string]string, 0, len(labelsMap))
			for k, v := range labelsMap {
				labelsArray = append(labelsArray, map[string]string{"label_name": k, "label_value": v.(string)})
			}
			component.UserConfig["resource_labels"] = labelsArray
		}
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func GcpCloudMonitoringDestinationToModel(plan *GcpCloudMonitoringDestinationModel, component *Destination) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.AckEnabled = NewBoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.CredentialsJSON = NewStringValue(component.UserConfig["credentials_json"].(string))
	plan.ProjectId = NewStringValue(component.UserConfig["project_id"].(string))
	plan.ResourceType = NewStringValue(component.UserConfig["resource_type"].(string))

	if component.UserConfig["resource_labels"] != nil {
		labelsArray, _ := component.UserConfig["resource_labels"].([]any)
		if len(labelsArray) > 0 {
			labelMap := make(map[string]any, len(labelsArray))
			for _, obj := range labelsArray {
				obj := obj.(map[string]any)
				key := obj["label_name"].(string)
				value := obj["label_value"].(string)
				labelMap[key] = value
			}
			labelType := GcpCloudMonitoringResourceSchema.Attributes["resource_labels"].GetType().(MapType).ElemType
			plan.ResourceLabels = NewMapValueMust(labelType, MapAnyToMapValues(labelMap))
		}
	}
}
