package destinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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

type ElasticSearchDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	AckEnabled   Bool   `tfsdk:"ack_enabled" user_config:"true"`
	Compression  String `tfsdk:"compression" user_config:"true"`
	Auth         Object `tfsdk:"auth" user_config:"true"`
	Endpoints    List   `tfsdk:"endpoints" user_config:"true"`
	Pipeline     String `tfsdk:"pipeline" user_config:"true"`
	Index        String `tfsdk:"index" user_config:"true"`
}

func ElasticSearchDestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents an ElasticSearch destination.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The compression strategy used on the encoded data prior to sending",
				Default:     stringdefault.StaticString("none"),
				Validators:  []validator.String{stringvalidator.OneOf("gzip", "none")},
			},
			"auth": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Configures ES authentication",
				Attributes: map[string]schema.Attribute{
					"strategy": schema.StringAttribute{
						Required:    true,
						Description: "The ES authentication strategy to use",
						Validators:  []validator.String{stringvalidator.OneOf("basic", "aws")},
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
					"access_key_id": schema.StringAttribute{
						Sensitive:   true,
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "The AWS access key id",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"secret_access_key": schema.StringAttribute{
						Sensitive:   true,
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "The AWS secret access key",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"region": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Description: "The AWS Region",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"endpoints": schema.ListAttribute{
				ElementType: StringType,
				Required:    true,
				Description: "An array of ElasticSearch endpoints",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"pipeline": schema.StringAttribute{
				Optional:    true,
				Description: "Name of an ES ingest pipeline to include",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"index": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Index to use when writing ES events (default = mezmo-%Y.%m.%d)",
				Default:     stringdefault.StaticString("mezmo-%Y.%m.%d"),
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		}, nil),
	}
}

func ElasticSearchDestinationFromModel(plan *ElasticSearchDestinationModel, previousState *ElasticSearchDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "elasticsearch",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"compression": plan.Compression.ValueString(),
				"index":       plan.Index.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
				"endpoints":   StringListValueToStringSlice(plan.Endpoints),
			},
		},
	}

	auth := MapValuesToMapAny(plan.Auth, &dd)
	component.UserConfig["auth"] = auth

	if auth["strategy"] == "basic" {
		if auth["user"] == "" || auth["password"] == "" {
			dd.AddError(
				"Error in plan",
				"Basic auth requires user and password fields to be defined")
		}
	} else {
		if auth["region"] == "" || auth["access_key_id"] == "" || auth["secret_access_key"] == "" {
			dd.AddError(
				"Error in plan",
				"AWS auth requires access_key_id, secret_access_key and region fields to be defined")
		}
	}

	if !plan.Pipeline.IsNull() {
		component.UserConfig["pipeline"] = plan.Pipeline.ValueString()
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func ElasticSearchDestinationToModel(plan *ElasticSearchDestinationModel, component *Destination) {
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

	plan.Compression = StringValue(component.UserConfig["compression"].(string))
	plan.Endpoints = SliceToStringListValue(component.UserConfig["endpoints"].([]any))
	plan.Index = StringValue(component.UserConfig["index"].(string))
	auth, _ := component.UserConfig["auth"].(map[string]any)
	if len(auth) > 0 {
		types := plan.Auth.AttributeTypes(context.Background())
		plan.Auth = basetypes.NewObjectValueMust(types, MapAnyToMapValues(auth))
	}

	if component.UserConfig["pipeline"] != nil {
		plan.Pipeline = StringValue(component.UserConfig["pipeline"].(string))
	}
}
