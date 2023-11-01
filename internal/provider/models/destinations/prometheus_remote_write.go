package destinations

import (
	"context"

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

const PROMETHEUS_REMOTE_WRITE_DESTINATION_TYPE_NAME = "prometheus_remote_write"
const PROMETHEUS_REMOTE_WRITE_DESTINATION_NODE_NAME = "prometheus-remote-write"

type PrometheusRemoteWriteDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled" user_config:"true"`
	Endpoint     String `tfsdk:"endpoint" user_config:"true"`
	Auth         Object `tfsdk:"auth" user_config:"true"`
}

var PrometheusRemoteWriteDestinationResourceSchema = schema.Schema{
	Description: "Represents Prometheus remote-write destination that publishes metrics to a " +
		"Prometheus endpoint",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"endpoint": schema.StringAttribute{
			Required: true,
			Description: "The full URI to make HTTP requests to. This should include the " +
				"protocol and host, but can also include the port, path, and any other valid " +
				"part of a URI. Example: http://example.org:8080/api/v1/push",
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"auth": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Configures authentication",
			Attributes: map[string]schema.Attribute{
				"strategy": schema.StringAttribute{
					Required:    true,
					Description: "The authentication strategy to use",
					Validators:  []validator.String{stringvalidator.OneOf("basic", "bearer")},
				},
				"user": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "The username for basic authentication",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"password": schema.StringAttribute{
					Sensitive:   true,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "The password to use for basic authentication",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"token": schema.StringAttribute{
					Sensitive:   true,
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "The token to use for bearer auth strategy",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
			},
		},
	}, nil),
}

func PrometheusRemoteWriteDestinationFromModel(plan *PrometheusRemoteWriteDestinationModel, previousState *PrometheusRemoteWriteDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        PROMETHEUS_REMOTE_WRITE_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"endpoint": plan.Endpoint.ValueString(),
			},
		},
	}

	if !plan.Auth.IsNull() {
		auth := MapValuesToMapAny(plan.Auth, &dd)
		component.UserConfig["auth"] = auth

		if auth["strategy"] == "basic" {
			if auth["user"] == "" || auth["password"] == "" {
				dd.AddError(
					"Error in plan",
					"Basic auth requires user and password fields to be defined")
			}
		} else {
			if auth["token"] == "" {
				dd.AddError(
					"Error in plan",
					"Bearer auth requires token field to be defined")
			}
		}
	} else {
		component.UserConfig["auth"] = map[string]string{"strategy": "none"}
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func PrometheusRemoteWriteDestinationToModel(plan *PrometheusRemoteWriteDestinationModel, component *Destination) {
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

	plan.Endpoint = StringValue(component.UserConfig["endpoint"].(string))
	auth, _ := component.UserConfig["auth"].(map[string]any)
	if len(auth) > 0 {
		if auth["strategy"] != "none" && auth["strategy"] != "" {
			attrTypes := plan.Auth.AttributeTypes(context.Background())
			authValues := MapAnyToMapValues(auth)
			// handle ConvertToTerraformModel calls
			if len(attrTypes) == 0 {
				attrTypes = PrometheusRemoteWriteDestinationResourceSchema.Attributes["auth"].GetType().(basetypes.ObjectType).AttrTypes
				authValues = MapAnyFillMissingValues(attrTypes, auth, MapKeys(attrTypes))
			}
			plan.Auth = basetypes.NewObjectValueMust(attrTypes, authValues)
		}
	}
}
