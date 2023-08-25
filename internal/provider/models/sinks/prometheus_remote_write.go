package sinks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type PrometheusRemoteWriteSinkModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled"`
	Endpoint     String `tfsdk:"endpoint"`
	Auth         Object `tfsdk:"auth"`
}

func PrometheusRemoteWriteSinkResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents Prometheus remote-write sink that publishes metrics to a " +
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
						Description: "The username for basic authentication",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"password": schema.StringAttribute{
						Sensitive:   true,
						Optional:    true,
						Description: "The password to use for basic authentication",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"token": schema.StringAttribute{
						Sensitive:   true,
						Optional:    true,
						Description: "The token to use for bearer auth strategy",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
		}, nil),
	}
}

func PrometheusRemoteWriteSinkFromModel(plan *PrometheusRemoteWriteSinkModel, previousState *PrometheusRemoteWriteSinkModel) (*Sink, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Sink{
		BaseNode: BaseNode{
			Type:        "prometheus-remote-write",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"endpoint": plan.Endpoint.ValueString(),
			},
		},
	}

	if !plan.Auth.IsNull() {
		auth, _ := MapValuesToMapStrings(plan.Auth, dd)
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

func PrometheusRemoteWriteSinkToModel(plan *PrometheusRemoteWriteSinkModel, component *Sink) {
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
	auth, _ := component.UserConfig["auth"].(map[string]string)
	if len(auth) > 0 {
		if auth["strategy"] != "none" && auth["strategy"] != "" {
			types := plan.Auth.AttributeTypes(context.Background())
			plan.Auth = basetypes.NewObjectValueMust(types, MapStringsToMapValues(auth))
		}
	}
}
