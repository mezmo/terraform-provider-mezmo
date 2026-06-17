package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v6/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/v6/internal/provider/models/modelutils"
)

const HTTP_CLIENT_SOURCE_TYPE_NAME = "http_client"
const HTTP_CLIENT_SOURCE_NODE_NAME = "http-client"

type HttpClientSourceModel struct {
	Id                 String `tfsdk:"id"`
	PipelineId         String `tfsdk:"pipeline_id"`
	Title              String `tfsdk:"title"`
	Description        String `tfsdk:"description"`
	GenerationId       Int64  `tfsdk:"generation_id"`
	Endpoint           String `tfsdk:"endpoint" user_config:"true"`
	Method             String `tfsdk:"method" user_config:"true"`
	ScrapeIntervalSecs Int64  `tfsdk:"scrape_interval_secs" user_config:"true"`
	ScrapeTimeoutSecs  Int64  `tfsdk:"scrape_timeout_secs" user_config:"true"`
	Auth               Object `tfsdk:"auth" user_config:"true"`
	Headers            Map    `tfsdk:"headers" user_config:"true"`
	Decoding           String `tfsdk:"decoding" user_config:"true"`
}

var HttpClientSourceResourceSchema = schema.Schema{
	Description: "Collect data by periodically polling an HTTP endpoint.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"endpoint": schema.StringAttribute{
			Required:    true,
			Description: "The URL to poll for data. Must be a valid URI.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				modelutils.URIValidator(),
			},
		},
		"method": schema.StringAttribute{
			Required:    true,
			Description: "The HTTP request method to use when polling the endpoint.",
			Validators: []validator.String{
				stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"),
			},
		},
		"scrape_interval_secs": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(30),
			Description: "The interval between polling requests, in seconds.",
			Validators: []validator.Int64{
				int64validator.Between(15, 86400),
			},
		},
		"scrape_timeout_secs": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(10),
			Description: "The timeout for each polling request, in seconds.",
			Validators: []validator.Int64{
				int64validator.Between(5, 60),
			},
		},
		"auth": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Configures HTTP authentication.",
			Attributes: map[string]schema.Attribute{
				"strategy": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("none"),
					Description: "The authentication strategy to use.",
					Validators: []validator.String{
						stringvalidator.OneOf("none", "basic", "bearer"),
					},
				},
				"user": schema.StringAttribute{
					Optional:    true,
					Description: "The basic authentication user.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"password": schema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "The basic authentication password.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"token": schema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "The bearer token.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
		"headers": schema.MapAttribute{
			Optional:    true,
			Description: "A key/value object describing a header name and its value.",
			ElementType: StringType,
			Validators: []validator.Map{
				mapvalidator.All(
					mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				),
			},
		},
		"decoding": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("bytes"),
			Description: "Configures how events are decoded from raw bytes.",
			Validators: []validator.String{
				stringvalidator.OneOf("bytes", "json"),
			},
		},
	}, nil),
}

func HttpClientSourceFromModel(plan *HttpClientSourceModel, previousState *HttpClientSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Source{
		BaseNode: BaseNode{
			Type:        HTTP_CLIENT_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"endpoint":             plan.Endpoint.ValueString(),
				"method":               plan.Method.ValueString(),
				"scrape_interval_secs": plan.ScrapeIntervalSecs.ValueInt64(),
				"scrape_timeout_secs":  plan.ScrapeTimeoutSecs.ValueInt64(),
				"decoding_codec":       plan.Decoding.ValueString(),
			},
		},
	}

	// Auth
	if !plan.Auth.IsNull() && !plan.Auth.IsUnknown() {
		auth := modelutils.MapValuesToMapAny(plan.Auth, &dd)
		component.UserConfig["auth"] = auth

		if auth["strategy"] == "basic" {
			_, hasUser := auth["user"].(string)
			_, hasPassword := auth["password"].(string)
			if !hasUser || !hasPassword {
				dd.AddError(
					"Error in plan",
					"Basic auth requires user and password fields to be defined")
			}
		} else if auth["strategy"] == "bearer" {
			_, hasToken := auth["token"].(string)
			if !hasToken {
				dd.AddError(
					"Error in plan",
					"Bearer auth requires token field to be defined")
			}
		}
	}

	// Headers
	if !plan.Headers.IsNull() {
		headerMap := modelutils.MapValuesToMapAny(plan.Headers, &dd)
		if len(headerMap) > 0 {
			headerArray := make([]map[string]string, 0, len(headerMap))
			for k, v := range headerMap {
				headerArray = append(headerArray, map[string]string{"header_name": k, "header_value": v.(string)})
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

func HttpClientSourceToModel(plan *HttpClientSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	plan.GenerationId = Int64Value(component.GenerationId)

	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	if component.UserConfig["endpoint"] != nil {
		plan.Endpoint = StringValue(component.UserConfig["endpoint"].(string))
	}
	if component.UserConfig["method"] != nil {
		plan.Method = StringValue(component.UserConfig["method"].(string))
	}
	if component.UserConfig["scrape_interval_secs"] != nil {
		plan.ScrapeIntervalSecs = Int64Value(int64(component.UserConfig["scrape_interval_secs"].(float64)))
	}
	if component.UserConfig["scrape_timeout_secs"] != nil {
		plan.ScrapeTimeoutSecs = Int64Value(int64(component.UserConfig["scrape_timeout_secs"].(float64)))
	}
	if component.UserConfig["decoding_codec"] != nil {
		plan.Decoding = StringValue(component.UserConfig["decoding_codec"].(string))
	}

	// Auth
	auth, _ := component.UserConfig["auth"].(map[string]any)
	if len(auth) > 0 {
		if auth["strategy"] != "none" && auth["strategy"] != "" {
			objT := plan.Auth.AttributeTypes(context.Background())
			if len(objT) == 0 {
				objT = HttpClientSourceResourceSchema.Attributes["auth"].GetType().(basetypes.ObjectType).AttrTypes
			}
			attributesToSetNil := []string{}
			switch auth["strategy"] {
			case "basic":
				attributesToSetNil = append(attributesToSetNil, "token")
			case "bearer":
				attributesToSetNil = append(attributesToSetNil, "user", "password")
			}

			mappedValues := modelutils.MapAnyToMapValues(auth)
			for _, a := range attributesToSetNil {
				mappedValues[a] = StringNull()
			}
			plan.Auth = basetypes.NewObjectValueMust(objT, mappedValues)
		}
	}

	// Headers
	if component.UserConfig["headers"] != nil {
		headerArray, _ := component.UserConfig["headers"].([]any)
		if len(headerArray) > 0 {
			headerMap := make(map[string]any, len(headerArray))
			for _, obj := range headerArray {
				obj := obj.(map[string]any)
				key := obj["header_name"].(string)
				value := obj["header_value"].(string)
				headerMap[key] = value
			}
			elType := plan.Headers.ElementType(context.Background())
			if elType == nil {
				elType = HttpClientSourceResourceSchema.Attributes["headers"].GetType().(basetypes.MapType).ElemType
			}
			plan.Headers = basetypes.NewMapValueMust(elType, modelutils.MapAnyToMapValues(headerMap))
		}
	}
}
