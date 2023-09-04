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

type HttpDestinationModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Inputs       List   `tfsdk:"inputs"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Uri          String `tfsdk:"uri"`
	Encoding     String `tfsdk:"encoding"`
	Compression  String `tfsdk:"compression"`
	Auth         Object `tfsdk:"auth"`
	Headers      Map    `tfsdk:"headers"`
	AckEnabled   Bool   `tfsdk:"ack_enabled"`
}

func HttpDestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents an HTTP destination.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Required: true,
				Description: "The full URI to make HTTP requests to. This should include the " +
					"protocol and host, but can also include the port, path, and any other valid " +
					"part of a URI.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"encoding": schema.StringAttribute{
				Optional:    true,
				Description: "The encoding to apply to the data",
				Computed:    true,
				Default:     stringdefault.StaticString("text"),
				Validators:  []validator.String{stringvalidator.OneOf("json", "ndjson", "text")},
			},
			"compression": schema.StringAttribute{
				Optional:    true,
				Description: "The compression strategy used on the encoded data prior to sending",
				Computed:    true,
				Default:     stringdefault.StaticString("none"),
				Validators:  []validator.String{stringvalidator.OneOf("gzip", "none")},
			},
			"auth": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Configures HTTP authentication",
				Attributes: map[string]schema.Attribute{
					"strategy": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.OneOf("basic", "bearer")},
					},
					"user": schema.StringAttribute{
						Optional:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"password": schema.StringAttribute{
						Sensitive:  true,
						Optional:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"token": schema.StringAttribute{
						Sensitive:  true,
						Optional:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"headers": schema.MapAttribute{
				Optional:    true,
				Description: "A key/value object describing a header name and its value",
				ElementType: StringType,
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

func HttpDestinationFromModel(plan *HttpDestinationModel, previousState *HttpDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "http",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"uri":         plan.Uri.ValueString(),
				"encoding":    plan.Encoding.ValueString(),
				"compression": plan.Compression.ValueString(),
				"ack_enabled": plan.AckEnabled.ValueBool(),
			},
		},
	}

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(basetypes.StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}
	if !plan.Auth.IsNull() {
		auth, _ := modelutils.MapValuesToMapStrings(plan.Auth, dd)
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

	if !plan.Headers.IsNull() {
		headerMap, ok := modelutils.MapValuesToMapStrings(plan.Headers, dd)
		if ok {
			headerArray := make([]map[string]string, 0, len(headerMap))
			for k, v := range headerMap {
				headerArray = append(headerArray, map[string]string{"header_name": k, "header_value": v})
			}
			component.UserConfig["headers"] = headerArray
		}
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func HttpDestinationToModel(plan *HttpDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	plan.Uri = StringValue(component.UserConfig["uri"].(string))
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.Compression = StringValue(component.UserConfig["compression"].(string))
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.GenerationId = Int64Value(component.GenerationId)

	auth, _ := component.UserConfig["auth"].(map[string]string)
	if len(auth) > 0 {
		if auth["strategy"] != "none" && auth["strategy"] != "" {
			types := plan.Auth.AttributeTypes(context.Background())
			plan.Auth = basetypes.NewObjectValueMust(types, modelutils.MapStringsToMapValues(auth))
		}
	}

	if component.UserConfig["headers"] != nil {
		headerArray, _ := component.UserConfig["headers"].([]map[string]string)
		if len(headerArray) > 0 {
			headerMap := make(map[string]string, len(headerArray))
			for _, obj := range headerArray {
				headerMap[obj["header_name"]] = obj["header_value"]
			}
			plan.Headers = basetypes.NewMapValueMust(MapType{}, modelutils.MapStringsToMapValues(headerMap))
		}
	}
}
