package destinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const HTTP_DESTINATION_TYPE_NAME = "http"
const HTTP_DESTINATION_NODE_NAME = HTTP_DESTINATION_TYPE_NAME

type HttpDestinationModel struct {
	Id            StringValue `tfsdk:"id"`
	PipelineId    StringValue `tfsdk:"pipeline_id"`
	Title         StringValue `tfsdk:"title"`
	Description   StringValue `tfsdk:"description"`
	Inputs        ListValue   `tfsdk:"inputs"`
	GenerationId  Int64Value  `tfsdk:"generation_id"`
	Uri           StringValue `tfsdk:"uri" user_config:"true"`
	Encoding      StringValue `tfsdk:"encoding" user_config:"true"`
	Compression   StringValue `tfsdk:"compression" user_config:"true"`
	Auth          ObjectValue `tfsdk:"auth" user_config:"true"`
	Headers       MapValue    `tfsdk:"headers" user_config:"true"`
	AckEnabled    BoolValue   `tfsdk:"ack_enabled" user_config:"true"`
	MaxBytes      Int64Value  `tfsdk:"max_bytes" user_config:"true"`
	TimeoutSecs   Int64Value  `tfsdk:"timeout_secs" user_config:"true"`
	Method        StringValue `tfsdk:"method" user_config:"true"`
	PayloadPrefix StringValue `tfsdk:"payload_prefix" user_config:"true"`
	PayloadSuffix StringValue `tfsdk:"payload_suffix" user_config:"true"`
	TLSProtocols  ListValue   `tfsdk:"tls_protocols" user_config:"true"`
	Proxy         ObjectValue `tfsdk:"proxy" user_config:"true"`
	RateLimiting  ObjectValue `tfsdk:"rate_limiting" user_config:"true"`
}

var HttpDestinationResourceSchema = schema.Schema{
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
					Computed:   true,
					Default:    stringdefault.StaticString(""),
					Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"password": schema.StringAttribute{
					Sensitive:  true,
					Optional:   true,
					Computed:   true,
					Default:    stringdefault.StaticString(""),
					Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"token": schema.StringAttribute{
					Sensitive:  true,
					Optional:   true,
					Computed:   true,
					Default:    stringdefault.StaticString(""),
					Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
				},
			},
		},
		"headers": schema.MapAttribute{
			Optional:    true,
			Description: "A key/value object describing a header name and its value",
			ElementType: StringType{},
			Validators: []validator.Map{
				mapvalidator.All(
					mapvalidator.KeysAre(stringvalidator.LengthAtLeast(1)),
					mapvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				),
			},
		},
		"max_bytes": schema.Int64Attribute{
			Optional:    true,
			Description: "The maximum number of uncompressed bytes when batching data to the destination",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.AtMost(1024 * 1024 * 2),
			},
		},
		"timeout_secs": schema.Int64Attribute{
			Optional:    true,
			Description: "The number of seconds before a destination write timeout.",
			Validators: []validator.Int64{
				int64validator.AtLeast(5),
			},
		},
		"method": schema.StringAttribute{
			Optional:    true,
			Description: "The HTTP method to use for the destination.",
			Validators:  []validator.String{stringvalidator.OneOf("post", "put", "patch", "delete", "get", "head", "options", "trace")},
		},
		"payload_prefix": schema.StringAttribute{
			Optional:    true,
			Description: "Add a prefix to the payload. Only used for serialized JSON chunks. This option also requires 'Payload Suffix' to form a valid JSON string.",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"payload_suffix": schema.StringAttribute{
			Optional:    true,
			Description: "Used in combination with 'Payload Prefix' to form valid JSON from the payload.",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"tls_protocols": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Description: "A list of ALPN protocols to use during TLS negotiation. They are attempted in the order they appear.",
			Validators: []validator.List{
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"proxy": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Proxy Settings",
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Turns Proxying on/off.",
				},
				"endpoint": schema.StringAttribute{
					Optional:    true,
					Description: "HTTP or HTTPS Endpoint to use for traffic.",
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"hosts_bypass_proxy": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Description: "A list of hosts to bypass proxying. Can be specified as a " +
						"domain name, IP address, or CIDR block. Wildcards are supported as " +
						"a dot (.) in domain names, or as a star (*) to match all hosts.",
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
					},
				},
			},
		},
		"rate_limiting": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Settings for controlling rate limiting to the destination.",
			Attributes: map[string]schema.Attribute{
				"request_limit": schema.Int64Attribute{
					Optional:    true,
					Description: "The max number of requests allowed within the specified.",
					Validators: []validator.Int64{
						int64validator.AtLeast(1),
					},
				},
				"duration_secs": schema.Int64Attribute{
					Optional:    true,
					Description: "The window of time used to apply 'Request Limit.",
					Validators: []validator.Int64{
						int64validator.AtLeast(1),
					},
				},
			},
		},
	}, nil),
}

func HttpDestinationFromModel(plan *HttpDestinationModel, previousState *HttpDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        HTTP_DESTINATION_NODE_NAME,
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
	user_config := component.UserConfig

	if !plan.Inputs.IsUnknown() {
		inputs := make([]string, 0)
		for _, v := range plan.Inputs.Elements() {
			value, _ := v.(StringValue)
			inputs = append(inputs, value.ValueString())
		}
		component.Inputs = inputs
	}
	if !plan.Auth.IsNull() {
		auth := MapValuesToMapAny(plan.Auth, &dd)
		user_config["auth"] = auth

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
		user_config["auth"] = map[string]string{"strategy": "none"}
	}

	if !plan.Headers.IsNull() {
		headerMap := MapValuesToMapAny(plan.Headers, &dd)
		if len(headerMap) > 0 {
			headerArray := make([]map[string]string, 0, len(headerMap))
			for k, v := range headerMap {
				headerArray = append(headerArray, map[string]string{"header_name": k, "header_value": v.(string)})
			}

			user_config["headers"] = headerArray
		}
	}
	advancedConfig := make(map[string]any)
	if !plan.MaxBytes.IsNull() && !plan.MaxBytes.IsUnknown() {
		advancedConfig["max_bytes"] = plan.MaxBytes.ValueInt64()
	}
	if !plan.MaxBytes.IsNull() && !plan.MaxBytes.IsUnknown() {
		advancedConfig["max_bytes"] = plan.MaxBytes.ValueInt64()
	}
	if !plan.TimeoutSecs.IsNull() && !plan.TimeoutSecs.IsUnknown() {
		advancedConfig["timeout_secs"] = plan.TimeoutSecs.ValueInt64()
	}
	if !plan.Method.IsNull() && !plan.Method.IsUnknown() {
		advancedConfig["method"] = plan.Method.ValueString()
	}
	if !plan.PayloadPrefix.IsNull() && !plan.PayloadPrefix.IsUnknown() {
		if plan.PayloadSuffix.IsNull() || plan.PayloadSuffix.IsUnknown() {
			dd.AddError(
				"Error in plan",
				"If 'payload_prefix' is set, 'payload_suffix' must be as well.",
			)
		} else {
			advancedConfig["payload_prefix"] = plan.PayloadPrefix.ValueString()
		}
	}
	if !plan.PayloadSuffix.IsNull() && !plan.PayloadSuffix.IsUnknown() {
		if plan.PayloadPrefix.IsNull() || plan.PayloadPrefix.IsUnknown() {
			dd.AddError(
				"Error in plan",
				"If 'payload_suffix' is set, 'payload_prefix' must be as well.",
			)
		} else {
			advancedConfig["payload_suffix"] = plan.PayloadSuffix.ValueString()
		}
	}
	if !plan.TLSProtocols.IsNull() && len(plan.TLSProtocols.Elements()) > 0 {
		advancedConfig["tls_protocols"] = StringListValueToStringSlice(plan.TLSProtocols)
	}
	if !plan.Proxy.IsNull() {
		proxy := MapValuesToMapAny(plan.Proxy, &dd)
		advancedConfig["proxy"] = proxy
	}
	if !plan.RateLimiting.IsNull() {
		rateLimiting := MapValuesToMapAny(plan.RateLimiting, &dd)
		advancedConfig["rate_limiting"] = rateLimiting
	}
	user_config["advanced_options"] = advancedConfig

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func HttpDestinationToModel(plan *HttpDestinationModel, component *Destination) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}

	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Uri = NewStringValue(component.UserConfig["uri"].(string))
	plan.Encoding = NewStringValue(component.UserConfig["encoding"].(string))
	plan.Compression = NewStringValue(component.UserConfig["compression"].(string))
	plan.AckEnabled = NewBoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.GenerationId = NewInt64Value(component.GenerationId)

	auth, _ := component.UserConfig["auth"].(map[string]any)
	if len(auth) > 0 {
		if auth["strategy"] != "none" && auth["strategy"] != "" {
			objT := plan.Auth.AttributeTypes(context.Background())
			if len(objT) == 0 {
				objT = HttpDestinationResourceSchema.Attributes["auth"].GetType().(ObjectType).AttrTypes
			}
			plan.Auth = NewObjectValueMust(objT, MapAnyToMapValues(auth))
		}
	}

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
			elType := plan.Headers.ElementType(nil)
			if elType == nil {
				elType = HttpDestinationResourceSchema.Attributes["headers"].GetType().(MapType).ElemType
			}
			plan.Headers = NewMapValueMust(elType, MapAnyToMapValues(headerMap))
		}
	}

	advanced, _ := component.UserConfig["advanced_options"].(map[string]any)
	if len(advanced) > 0 {
		if advanced["max_bytes"] != nil {
			plan.MaxBytes = NewInt64Value(int64(advanced["max_bytes"].(float64)))
		}
		if advanced["timeout_secs"] != nil {
			plan.TimeoutSecs = NewInt64Value(int64(advanced["timeout_secs"].(float64)))
		}
		if advanced["method"] != nil {
			plan.Method = NewStringValue(advanced["method"].(string))
		}
		if advanced["payload_prefix"] != nil {
			plan.PayloadPrefix = NewStringValue(advanced["payload_prefix"].(string))
		}
		if advanced["payload_suffix"] != nil {
			plan.PayloadSuffix = NewStringValue(advanced["payload_suffix"].(string))
		}
		if advanced["tls_protocols"] != nil {
			plan.TLSProtocols = SliceToStringListValue(advanced["tls_protocols"].([]any))
		}

		if advanced["proxy"] != nil {
			component_map, _ := advanced["proxy"].(map[string]any)
			plan_map := map[string]attr.Value{
				"enabled":            types.BoolNull(),
				"endpoint":           types.StringNull(),
				"hosts_bypass_proxy": types.ListNull(types.StringType),
			}

			if component_map["enabled"] != nil {
				plan_map["enabled"] = types.BoolValue(component_map["enabled"].(bool))
			}
			if component_map["endpoint"] != nil {
				plan_map["endpoint"] = types.StringValue(component_map["endpoint"].(string))
			}

			if component_map["hosts_bypass_proxy"] != nil {
				list := make([]attr.Value, 0)
				for _, v := range component_map["hosts_bypass_proxy"].([]any) {
					value, _ := v.(string)
					list = append(list, types.StringValue(value))
				}
				plan_map["hosts_bypass_proxy"] = types.ListValueMust(types.StringType, list)
			}

			attrTypes := plan.Proxy.AttributeTypes(context.Background())
			if len(attrTypes) == 0 {
				attrTypes = HttpDestinationResourceSchema.Attributes["proxy"].GetType().(ObjectType).AttrTypes
			}
			plan.Proxy = NewObjectValueMust(attrTypes, plan_map)
		}

		if advanced["rate_limiting"] != nil {
			component_map, _ := advanced["rate_limiting"].(map[string]any)
			plan_map := map[string]attr.Value{
				"request_limit": types.Int64Null(),
				"duration_secs": types.Int64Null(),
			}

			if component_map["request_limit"] != nil {
				plan_map["request_limit"] = types.Int64Value(int64(component_map["request_limit"].(float64)))
			}
			if component_map["duration_secs"] != nil {
				plan_map["duration_secs"] = types.Int64Value(int64(component_map["duration_secs"].(float64)))
			}
			attrTypes := plan.RateLimiting.AttributeTypes(context.Background())
			if len(attrTypes) == 0 {
				attrTypes = HttpDestinationResourceSchema.Attributes["rate_limiting"].GetType().(ObjectType).AttrTypes
			}
			plan.RateLimiting = NewObjectValueMust(attrTypes, plan_map)
		}
	}
}
