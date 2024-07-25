package destinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const GCP_CLOUD_OPERATIONS_DESTINATION_TYPE_NAME = "gcp_cloud_operations"
const GCP_CLOUD_OPERATIONS_DESTINATION_NODE_NAME = "gcp-cloud-operations"

type GcpCloudOperationsDestinationModel struct {
	Id              String `tfsdk:"id"`
	PipelineId      String `tfsdk:"pipeline_id"`
	Title           String `tfsdk:"title"`
	Description     String `tfsdk:"description"`
	Inputs          List   `tfsdk:"inputs"`
	GenerationId    Int64  `tfsdk:"generation_id"`
	AckEnabled      Bool   `tfsdk:"ack_enabled" user_config:"true"`
	CredentialsJSON String `tfsdk:"credentials_json" user_config:"true"`
	LogId           String `tfsdk:"log_id" user_config:"true"`
	ProjectId       String `tfsdk:"project_id" user_config:"true"`
	ResourceType    String `tfsdk:"resource_type" user_config:"true"`
	ResourceLabels  Map    `tfsdk:"resource_labels" user_config:"true"`
}

var GcpCloudOperationsResourceSchema = schema.Schema{
	Description: "Publish log events to GCP Cloud Operations",
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
		"log_id": schema.StringAttribute{
			Required:    true,
			Description: "Concise reference for the log stream name.",
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
			ElementType: StringType,
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

func GcpCloudOperationsDestinationFromModel(plan *GcpCloudOperationsDestinationModel, previousState *GcpCloudOperationsDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        GCP_CLOUD_OPERATIONS_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":      plan.AckEnabled.ValueBool(),
				"credentials_json": plan.CredentialsJSON.ValueString(),
				"project_id":       plan.ProjectId.ValueString(),
				"log_id":           plan.LogId.ValueString(),
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

func GcpCloudOperationsDestinationToModel(plan *GcpCloudOperationsDestinationModel, component *Destination) {
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
	plan.CredentialsJSON = StringValue(component.UserConfig["credentials_json"].(string))
	plan.ProjectId = StringValue(component.UserConfig["project_id"].(string))
	plan.LogId = StringValue(component.UserConfig["log_id"].(string))
	plan.ResourceType = StringValue(component.UserConfig["resource_type"].(string))

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
			labelType := plan.ResourceLabels.ElementType(context.Background())
			if labelType == nil {
				labelType = GcpCloudMonitoringResourceSchema.Attributes["resource_labels"].GetType().(basetypes.MapType).ElemType
			}
			plan.ResourceLabels = basetypes.NewMapValueMust(labelType, MapAnyToMapValues(labelMap))
		}
	}
}
