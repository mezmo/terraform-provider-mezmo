package sinks

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type KafkaSinkModel struct {
	Id            String `tfsdk:"id"`
	PipelineId    String `tfsdk:"pipeline_id"`
	Title         String `tfsdk:"title"`
	Description   String `tfsdk:"description"`
	Inputs        List   `tfsdk:"inputs"`
	GenerationId  Int64  `tfsdk:"generation_id"`
	Encoding      String `tfsdk:"encoding"`
	Compression   String `tfsdk:"compression"`
	EventKeyField String `tfsdk:"event_key_field"`
	Brokers       List   `tfsdk:"brokers"`
	Topic         String `tfsdk:"topic"`
	TLSEnabled    Bool   `tfsdk:"tls_enabled"`
	SASL          Object `tfsdk:"sasl"`
	AckEnabled    Bool   `tfsdk:"ack_enabled"`
}

func KafkaSinkResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents a Kafka sink.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"encoding": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("text"),
				Description: "The encoding to apply to the data.",
				Validators: []validator.String{
					stringvalidator.OneOf("json", "text"),
				},
			},
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("none"),
				Description: "The compression strategy used on the encoded data prior to sending.",
				Validators: []validator.String{
					stringvalidator.OneOf("gzip", "lz4", "snappy", "zstd", "none"),
				},
			},
			"event_key_field": schema.StringAttribute{
				Optional:    true,
				Description: "The field in the log whose value is used as Kafka's event key.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"brokers": schema.ListNestedAttribute{
				Required:    true,
				Description: "The Kafka brokers to connect to.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Required:    true,
							Description: "The host of the Kafka broker.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
								stringvalidator.LengthAtMost(255),
							},
						},
						"port": schema.Int64Attribute{
							Required:    true,
							Description: "The port of the Kafka broker.",
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"topic": schema.StringAttribute{
				Required:    true,
				Description: "The name of the topic to publish to.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"tls_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to use TLS when connecting to Kafka.",
			},
			"sasl": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The SASL configuration to use when connecting to Kafka.",
				Attributes: map[string]schema.Attribute{
					"mechanism": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The SASL mechanism to use when connecting to Kafka.",
						Validators: []validator.String{
							stringvalidator.OneOf("PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512"),
						},
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The SASL username to use when connecting to Kafka.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The SASL password to use when connecting to Kafka.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
		}, nil),
	}
}

func KafkaSinkFromModel(plan *KafkaSinkModel, previousState *KafkaSinkModel) (*Sink, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Sink{
		BaseNode: BaseNode{
			Type:        "kafka",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"encoding":    plan.Encoding.ValueString(),
				"compression": plan.Compression.ValueString(),
				"topic":       plan.Topic.ValueString(),
				"tls_enabled": plan.TLSEnabled.ValueBool(),
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

	if !plan.EventKeyField.IsNull() {
		component.UserConfig["event_key_field"] = plan.EventKeyField.ValueString()
	}

	brokers, dd := modelutils.BrokersFromModelList(plan.Brokers, dd)
	component.UserConfig["brokers"] = brokers

	sasl := plan.SASL.Attributes()
	if len(sasl) > 0 {
		component.UserConfig["sasl_enabled"] = true
		component.UserConfig["sasl_username"] = modelutils.GetAttributeValue[String](sasl, "username").ValueString()
		component.UserConfig["sasl_password"] = modelutils.GetAttributeValue[String](sasl, "password").ValueString()
		component.UserConfig["sasl_mechanism"] = modelutils.GetAttributeValue[String](sasl, "mechanism").ValueString()
	} else {
		component.UserConfig["sasl_enabled"] = false
	}

	if previousState != nil {
		// Set generated fields
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}
	return &component, dd
}

func KafkaSinkToModel(plan *KafkaSinkModel, component *Sink) {
	plan.Id = StringValue(component.Id)
	plan.GenerationId = Int64Value(component.GenerationId)

	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	plan.Inputs = modelutils.SliceToStringListValue(component.Inputs)
	plan.Brokers = modelutils.BrokersToModelList(plan.Brokers.ElementType(context.Background()), component.UserConfig["brokers"].([]interface{}))
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.Compression = StringValue(component.UserConfig["compression"].(string))
	plan.Topic = StringValue(component.UserConfig["topic"].(string))
	plan.TLSEnabled = BoolValue(component.UserConfig["tls_enabled"].(bool))
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))

	if component.UserConfig["event_key_field"] != nil {
		plan.EventKeyField = StringValue(component.UserConfig["event_key_field"].(string))
	}

	if component.UserConfig["sasl_enabled"] != nil {
		sasl_enabled, _ := component.UserConfig["sasl_enabled"].(bool)
		if sasl_enabled {
			plan.SASL = modelutils.KafkaSinkSASLToModel(plan.SASL.AttributeTypes(context.Background()), component.UserConfig)
		}
	}
}
