package destinations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type SplunkHecLogsDestinationModel struct {
	Id                   String `tfsdk:"id"`
	PipelineId           String `tfsdk:"pipeline_id"`
	Title                String `tfsdk:"title"`
	Description          String `tfsdk:"description"`
	Inputs               List   `tfsdk:"inputs"`
	GenerationId         Int64  `tfsdk:"generation_id"`
	AckEnabled           Bool   `tfsdk:"ack_enabled" user_config:"true"`
	Compression          String `tfsdk:"compression" user_config:"true"`
	Endpoint             String `tfsdk:"endpoint" user_config:"true"`
	Token                String `tfsdk:"token" user_config:"true"`
	HostField            String `tfsdk:"host_field" user_config:"true"`
	TimestampField       String `tfsdk:"timestamp_field" user_config:"true"`
	TlsVerifyCertificate Bool   `tfsdk:"tls_verify_certificate" user_config:"true"`
	Source               Object `tfsdk:"source" user_config:"true"`
	SourceType           Object `tfsdk:"source_type" user_config:"true"`
	Index                Object `tfsdk:"index" user_config:"true"`
}

var splunkValueTypeAttributes = map[string]schema.Attribute{
	"field": schema.StringAttribute{
		Optional:    true,
		Description: "The field path to use",
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "The fixed value to use",
		Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
	},
}

func SplunkHecLogsDestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Publishes log events to a Splunk HTTP Event Collector",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The compression strategy used on the encoded data prior to sending",
				Default:     stringdefault.StaticString("none"),
				Validators:  []validator.String{stringvalidator.OneOf("gzip", "none")},
			},
			"endpoint": schema.StringAttribute{
				Required: true,
				Description: "The base URL for the Splunk instance. The collector path, such as " +
					"`/services/collector/events`, will be automatically inferred from the " +
					"destination's configuration.",
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"token": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The default token to authenticate to Splunk HEC",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"tls_verify_certificate": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Verify TLS Certificate",
				Default:     booldefault.StaticBool(true),
			},
			"host_field": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The field that contains the hostname to include in the event",
				Default:     stringdefault.StaticString("metadata.host"),
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"timestamp_field": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The field that contains the timestamp to include in the event",
				Default:     stringdefault.StaticString("metadata.time"),
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"source": schema.SingleNestedAttribute{
				Optional: true,
				Description: "The source of events sent to this destination. This is typically the filename " +
					"the logs originated from. Use the field path \"metadata.source\" to use the" +
					" upstream source value from a HEC log source",
				Attributes: splunkValueTypeAttributes,
			},
			"source_type": schema.SingleNestedAttribute{
				Optional: true,
				Description: "The sourcetype of events sent to this destination. Use the field path" +
					" \"metadata.sourcetype\" to use the upstream sourcetype value from a HEC" +
					" log source",
				Attributes: splunkValueTypeAttributes,
			},
			"index": schema.SingleNestedAttribute{
				Optional: true,
				Description: "The name of the index to send events to. Use the field path " +
					" \"metadata.index\" to use the upstream index value from a HEC log source",
				Attributes: splunkValueTypeAttributes,
			},
		}, nil),
	}
}

func SplunkHecLogsDestinationFromModel(plan *SplunkHecLogsDestinationModel, previousState *SplunkHecLogsDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "splunk-hec-logs",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":            plan.AckEnabled.ValueBool(),
				"compression":            plan.Compression.ValueString(),
				"url":                    plan.Endpoint.ValueString(),
				"token":                  plan.Token.ValueString(),
				"tls_verify_certificate": plan.TlsVerifyCertificate.ValueBool(),
				"host_field":             plan.HostField.ValueString(),
				"timestamp_field":        plan.TimestampField.ValueString(),
				"source":                 map[string]any{"value_type": "none"},
				"sourcetype":             map[string]any{"value_type": "none"},
				"index":                  map[string]any{"value_type": "none"},
			},
		},
	}

	splunkValueTypeFromModel(&plan.Source, &component, "source", "source", &dd)
	splunkValueTypeFromModel(&plan.SourceType, &component, "sourcetype", "source_type", &dd)
	splunkValueTypeFromModel(&plan.Index, &component, "index", "index", &dd)

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func splunkValueTypeFromModel(
	planValue *Object,
	component *Destination,
	configName string,
	attrName string,
	dd *diag.Diagnostics,
) {
	if !planValue.IsNull() {
		m := MapValuesToMapAny(*planValue, dd)
		target := component.UserConfig[configName].(map[string]any)
		if m["field"] != "" {
			target["value_type"] = "field"
			target["value"] = m["field"]
		} else if m["value"] != "" {
			target["value_type"] = "value"
			target["value"] = m["value"]
		} else {
			dd.AddError(
				"Error in plan",
				fmt.Sprintf("%s requires field or value to be defined", attrName))
		}
	}
}

func splunkValueTypeToModel(planObj *Object, component *Destination, configName string) Object {
	m, _ := component.UserConfig[configName].(map[string]any)
	types := planObj.AttributeTypes(context.Background())
	if len(m) > 0 && m["value_type"] != "none" {
		planMap := map[string]attr.Value{
			"field": StringNull(),
			"value": StringNull(),
		}
		if m["value_type"] == "field" {
			planMap["field"] = StringValue(m["value"].(string))
		} else if m["value_type"] == "value" {
			planMap["value"] = StringValue(m["value"].(string))
		}
		return basetypes.NewObjectValueMust(types, planMap)
	}

	return ObjectNull(types)
}

func SplunkHecLogsDestinationToModel(plan *SplunkHecLogsDestinationModel, component *Destination) {
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
	plan.Endpoint = StringValue(component.UserConfig["url"].(string))
	plan.Token = StringValue(component.UserConfig["token"].(string))
	plan.TlsVerifyCertificate = BoolValue(component.UserConfig["tls_verify_certificate"].(bool))
	plan.HostField = StringValue(component.UserConfig["host_field"].(string))
	plan.TimestampField = StringValue(component.UserConfig["timestamp_field"].(string))
	plan.Source = splunkValueTypeToModel(&plan.Source, component, "source")
	plan.SourceType = splunkValueTypeToModel(&plan.SourceType, component, "sourcetype")
	plan.Index = splunkValueTypeToModel(&plan.Index, component, "index")
}
