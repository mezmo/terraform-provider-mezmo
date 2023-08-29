package sources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type KafkaSourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Brokers      List   `tfsdk:"brokers"`
	Topics       List   `tfsdk:"topics"`
	GroupId      String `tfsdk:"group_id"`
	TLSEnabled   Bool   `tfsdk:"tls_enabled"`
	SASL         Object `tfsdk:"sasl"`
	Decoding     String `tfsdk:"decoding"`
}

func KafkaSourceResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Represents a Kafka source.",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
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
			"topics": schema.ListAttribute{
				Required:    true,
				Description: "The Kafka topics to consume from.",
				ElementType: StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 256),
					),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "The Kafka consumer group ID to use.",
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
			"decoding": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("json"),
				Description: "The decoding method for converting frames into data events.",
				Validators: []validator.String{
					stringvalidator.OneOf("bytes", "json"),
				},
			},
		}, nil),
	}
}

func KafkaSourceFromModel(plan *KafkaSourceModel, previousState *KafkaSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	brokers, dd := BrokersFromModelList(plan.Brokers, dd)

	topics := make([]string, 0, len(plan.Topics.Elements()))
	dd = plan.Topics.ElementsAs(context.Background(), &topics, false)

	component := Source{
		BaseNode: BaseNode{
			Type:        "kafka",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"brokers":        brokers,
				"topics":         topics,
				"group_id":       plan.GroupId.ValueString(),
				"tls_enabled":    plan.TLSEnabled.ValueBool(),
				"decoding_codec": plan.Decoding.ValueString(),
			},
		},
	}

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

func KafkaSourceToModel(plan *KafkaSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	plan.GenerationId = Int64Value(component.GenerationId)

	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	plan.Brokers = BrokersToModelList(plan, component.UserConfig["brokers"].([]interface{}))
	plan.Topics, _ = ListValueFrom(context.Background(), StringType, component.UserConfig["topics"])
	plan.GroupId = StringValue(component.UserConfig["group_id"].(string))
	plan.TLSEnabled = BoolValue(component.UserConfig["tls_enabled"].(bool))

	if component.UserConfig["sasl_enabled"] != nil {
		sasl_enabled, _ := component.UserConfig["sasl_enabled"].(bool)
		if sasl_enabled {
			types := plan.SASL.AttributeTypes(context.Background())
			sasl := map[string]attr.Value{}
			if component.UserConfig["sasl_username"] != nil {
				sasl["username"] = StringValue(component.UserConfig["sasl_username"].(string))
			}
			if component.UserConfig["sasl_password"] != nil {
				sasl["password"] = StringValue(component.UserConfig["sasl_password"].(string))
			}
			if component.UserConfig["sasl_mechanism"] != nil {
				sasl["mechanism"] = StringValue(component.UserConfig["sasl_mechanism"].(string))
			}
			plan.SASL = basetypes.NewObjectValueMust(types, sasl)
		}
	}

	plan.Decoding = StringValue(component.UserConfig["decoding_codec"].(string))
}

func BrokersFromModelList(Brokers List, dd diag.Diagnostics) ([]map[string]any, diag.Diagnostics) {
	output := make([]map[string]any, 0)
	elements := Brokers.Elements()
	for _, b := range elements {
		broker := map[string]any{}
		attrs := b.(basetypes.ObjectValue).Attributes()
		for k, v := range attrs {
			switch v.(type) {
			case String:
				value, ok := attrs[k].(basetypes.StringValue)
				if !ok {
					dd.AddError(
						"Could not look up attribute value",
						fmt.Sprintf("Cannot cast key %s to a string value. Please report this to Mezmo.", k),
					)
					continue
				}
				broker[k] = value.ValueString()
			case Int64:
				value, ok := attrs[k].(basetypes.Int64Value)
				if !ok {
					dd.AddError(
						"Could not look up attribute value",
						fmt.Sprintf("Cannot cast key %s to an int value. Please report this to Mezmo.", k),
					)
					continue
				}
				broker[k] = value.ValueInt64()
			}
		}
		output = append(output, broker)
	}
	return output, dd
}

func BrokersToModelList(plan *KafkaSourceModel, brokers []interface{}) List {
	output := make([]attr.Value, 0)
	for _, v := range brokers {
		broker_raw := v.(map[string]interface{})
		broker_map := map[string]attr.Value{
			"host": StringValue(broker_raw["host"].(string)),
			"port": Int64Value(int64(broker_raw["port"].(float64))),
		}
		broker := basetypes.NewObjectValueMust(
			map[string]attr.Type{"host": StringType, "port": Int64Type},
			broker_map)
		output = append(output, broker)
	}
	brokersList, _ := ListValue(plan.Brokers.ElementType(context.Background()), output)
	return brokersList
}
