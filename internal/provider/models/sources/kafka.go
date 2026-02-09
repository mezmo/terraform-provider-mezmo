package sources

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
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/modelutils"
)

const KAFKA_SOURCE_TYPE_NAME = "kafka"
const KAFKA_SOURCE_NODE_NAME = KAFKA_SOURCE_TYPE_NAME

type KafkaSourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Brokers      List   `tfsdk:"brokers" user_config:"true"`
	Topics       List   `tfsdk:"topics" user_config:"true"`
	GroupId      String `tfsdk:"group_id" user_config:"true"`
	TLSEnabled   Bool   `tfsdk:"tls_enabled" user_config:"true"`
	SASL         Object `tfsdk:"sasl" user_config:"true"`
	Decoding     String `tfsdk:"decoding" user_config:"true"`
}

var KafkaSourceResourceSchema = schema.Schema{
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

func KafkaSourceFromModel(plan *KafkaSourceModel, previousState *KafkaSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	brokers, dd := modelutils.BrokersFromModelList(plan.Brokers, dd)

	topics := make([]string, 0, len(plan.Topics.Elements()))
	dd = plan.Topics.ElementsAs(context.Background(), &topics, false)

	component := Source{
		BaseNode: BaseNode{
			Type:        KAFKA_SOURCE_NODE_NAME,
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

	elemType := plan.Brokers.ElementType(context.Background())
	if elemType == nil {
		elemType = KafkaSourceResourceSchema.Attributes["brokers"].GetType().(basetypes.ListType).ElemType
	}
	plan.Brokers = modelutils.BrokersToModelList(elemType, component.UserConfig["brokers"].([]interface{}))
	plan.Topics, _ = ListValueFrom(context.Background(), StringType, component.UserConfig["topics"])
	plan.GroupId = StringValue(component.UserConfig["group_id"].(string))
	plan.TLSEnabled = BoolValue(component.UserConfig["tls_enabled"].(bool))

	if component.UserConfig["sasl_enabled"] != nil {
		sasl_enabled, _ := component.UserConfig["sasl_enabled"].(bool)
		if sasl_enabled {
			planType := plan.SASL.AttributeTypes(context.Background())
			if len(planType) == 0 {
				planType = KafkaSourceResourceSchema.Attributes["sasl"].GetType().(basetypes.ObjectType).AttributeTypes()
			}
			plan.SASL = modelutils.KafkaDestinationSASLToModel(planType, component.UserConfig)
		}
	}

	plan.Decoding = StringValue(component.UserConfig["decoding_codec"].(string))
}
