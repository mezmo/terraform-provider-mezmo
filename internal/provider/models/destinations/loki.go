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
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type LokiDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Endpoint     String `tfsdk:"endpoint"`
	Path         String `tfsdk:"path"`
	Encoding     String `tfsdk:"encoding"`
	Auth         Object `tfsdk:"auth"`
	Labels       Map    `tfsdk:"labels"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled"`
}

func LokiDestinationResourceSchema() schema.Schema {
	return schema.Schema{
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
}

func LokiFromModel(plan *LokiDestinationModel, previousState *LokiDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Destination{
		BaseNode: BaseNode{
			Type:        "loki",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      modelutils.StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"endpoint":    plan.Endpoint.ValueString(),
				"path":        plan.Path.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
			},
		},
	}

	if !plan.Auth.IsNull() {
		component.UserConfig["auth"], _ = modelutils.MapValuesToMapStrings(plan.Auth, dd)
	}

	if !plan.Encoding.IsNull() {
		component.UserConfig["encoding"] = map[string]string{"codec": plan.Encoding.ValueString()}
	}

	if !plan.Labels.IsNull() {
		lablesMap, ok := modelutils.MapValuesToMapStrings(plan.Labels, dd)
		if ok == false {
			labelsArray := make([]map[string]string, 0, len(lablesMap))
			for k, v := range lablesMap {
				labelsArray = append(labelsArray, map[string]string{"label_name": k, "label_value": v})
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
	plan.Title = StringValue(component.Title)
	plan.Description = StringValue(component.Description)
	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.Endpoint = StringValue(component.UserConfig["endpoint"].(string))
	plan.Path = StringValue(component.UserConfig["path"].(string))
	if component.UserConfig["auth"] != nil {
		values, _ := component.UserConfig["auth"].(map[string]string)
		if len(values) > 0 {
			types := plan.Auth.AttributeTypes(context.Background())
			plan.Auth = basetypes.NewObjectValueMust(types, modelutils.MapStringsToMapValues(values))
		}
	}
	if component.UserConfig["encoding"] != nil {
		codecValue := component.UserConfig["encoding"].(map[string]interface{})
		plan.Encoding = StringValue(codecValue["codec"].(string))
	}
	if component.UserConfig["labels"] != nil {
		labelsArray, _ := component.UserConfig["labels"].([]map[string]string)
		if len(labelsArray) > 0 {
			labelMap := make(map[string]string, len(labelsArray))
			for _, obj := range labelsArray {
				labelMap[obj["label_name"]] = obj["label_value"]
			}
			plan.Labels = basetypes.NewMapValueMust(MapType{}, modelutils.MapStringsToMapValues(labelMap))
		}
	}
}
