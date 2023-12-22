package destinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const LOKI_DESTINATION_TYPE_NAME = "loki"
const LOKI_DESTINATION_NODE_NAME = "loki"

type LokiDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Endpoint     String `tfsdk:"endpoint" user_config:"true"`
	Path         String `tfsdk:"path" user_config:"true"`
	Encoding     String `tfsdk:"encoding" user_config:"true"`
	Auth         Object `tfsdk:"auth" user_config:"true"`
	Labels       Map    `tfsdk:"labels" user_config:"true"`
	Inputs       List   `tfsdk:"inputs" user_config:"true"`
	AckEnabled   Bool   `tfsdk:"ack_enabled" user_config:"true"`
}

var LokiDestinationResourceSchema = schema.Schema{
	Description: "Publish log events to Loki",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"auth": schema.SingleNestedAttribute{
			Required:    true,
			Description: "The authentication strategy to use (only basic supported)",
			Attributes: map[string]schema.Attribute{
				"strategy": schema.StringAttribute{
					Required:    true,
					Description: "The authentication strategy to use (only basic supported)",
					Validators:  []validator.String{stringvalidator.OneOf("basic")},
				},
				"user": schema.StringAttribute{
					Required:    true,
					Description: "The basic authentication user",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"password": schema.StringAttribute{
					Sensitive:   true,
					Required:    true,
					Description: "The basic authentication password",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
			},
		},
		"endpoint": schema.StringAttribute{
			Required:    true,
			Description: "The Loki base URL",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"encoding": schema.StringAttribute{
			Required:    true,
			Description: "Configures how event are encoded",
			Validators:  []validator.String{stringvalidator.OneOf("json", "text")},
		},
		"path": schema.StringAttribute{
			Optional:    true,
			Description: "The path appended to the Loki base URL, (defaults to /loki/api/v1/push)",
			Computed:    true,
			Default:     stringdefault.StaticString("/loki/api/v1/push"),
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"labels": schema.MapAttribute{
			Required:    true,
			ElementType: StringType,
			Description: "Key/Value pair used to describe Loki data",
			Validators: []validator.Map{
				mapvalidator.All(
					mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				),
			},
		},
	}, nil),
}

func LokiFromModel(plan *LokiDestinationModel, previousState *LokiDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Destination{
		BaseNode: BaseNode{
			Type:        LOKI_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"endpoint":    plan.Endpoint.ValueString(),
				"path":        plan.Path.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
			},
		},
	}

	if !plan.Auth.IsNull() {
		component.UserConfig["auth"] = MapValuesToMapAny(plan.Auth, &dd)
	}

	if !plan.Encoding.IsNull() {
		component.UserConfig["encoding"] = map[string]string{"codec": plan.Encoding.ValueString()}
	}

	if !plan.Labels.IsNull() {
		lablesMap := MapValuesToMapAny(plan.Labels, &dd)
		if !dd.HasError() {
			labelsArray := make([]map[string]string, 0, len(lablesMap))
			for k, v := range lablesMap {
				labelsArray = append(labelsArray, map[string]string{"label_name": k, "label_value": v.(string)})
			}
			component.UserConfig["labels"] = labelsArray
		}
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func LokiDestinationToModel(plan *LokiDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.Endpoint = StringValue(component.UserConfig["endpoint"].(string))
	plan.Path = StringValue(component.UserConfig["path"].(string))
	if component.UserConfig["auth"] != nil {
		values, _ := component.UserConfig["auth"].(map[string]any)
		if len(values) > 0 {
			authAttrTypes := plan.Auth.AttributeTypes(context.Background())
			if len(authAttrTypes) == 0 {
				authAttrTypes = LokiDestinationResourceSchema.Attributes["auth"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.Auth = basetypes.NewObjectValueMust(
				authAttrTypes,
				MapAnyFillMissingValues(authAttrTypes, values, MapKeys(authAttrTypes)),
			)
		}
	}
	if component.UserConfig["encoding"] != nil {
		codecValue := component.UserConfig["encoding"].(map[string]interface{})
		plan.Encoding = StringValue(codecValue["codec"].(string))
	}
	if component.UserConfig["labels"] != nil {
		labelsArray, _ := component.UserConfig["labels"].([]any)
		if len(labelsArray) > 0 {
			labelMap := make(map[string]any, len(labelsArray))
			for _, obj := range labelsArray {
				obj := obj.(map[string]any)
				key := obj["label_name"].(string)
				value := obj["label_value"].(string)
				labelMap[key] = value
			}
			labelType := plan.Labels.ElementType(context.Background())
			if labelType == nil {
				labelType = LokiDestinationResourceSchema.Attributes["labels"].GetType().(basetypes.MapType).ElemType
			}
			plan.Labels = basetypes.NewMapValueMust(labelType, MapAnyToMapValues(labelMap))
		}
	}
}
