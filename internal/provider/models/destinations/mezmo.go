package destinations

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type MezmoDestinationModel struct {
	Id                    String `tfsdk:"id"`
	PipelineId            String `tfsdk:"pipeline_id"`
	Title                 String `tfsdk:"title"`
	Description           String `tfsdk:"description"`
	Inputs                List   `tfsdk:"inputs"`
	GenerationId          Int64  `tfsdk:"generation_id"`
	AckEnabled            Bool   `tfsdk:"ack_enabled" user_config:"true"`
	Host                  String `tfsdk:"host" user_config:"true"`
	IngestionKey          String `tfsdk:"ingestion_key" user_config:"true"`
	Query                 Object `tfsdk:"query" user_config:"true"`
	LogConstructionScheme String `tfsdk:"log_construction_scheme" user_config:"true"`
	ExplicitSchemeOptions Object `tfsdk:"explicit_scheme_options" user_config:"true"`
}

var log_construction_schemes = map[string]string{
	"explicit":     "Explicit field selection",
	"pass-through": "Message pass-through",
}

var MezmoDestinationResourceSchema = schema.Schema{
	Description: "Represents a Mezmo destination.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"host": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("logs.logdna.com"),
			Description: "The host for your Log Analysis environment",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"ingestion_key": schema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "Ingestion key",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"query": schema.SingleNestedAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Query Parameters",
			Attributes: map[string]schema.Attribute{
				"hostname": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "Hostname string or template to attach to logs",
					Default:     stringdefault.StaticString("mezmo"),
					Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
				},
				"tags": schema.ListAttribute{
					ElementType: StringType,
					Optional:    true,
					Description: "List of tag strings or templates to attach to logs",
					Validators: []validator.List{
						listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
					},
				},
				"ip": schema.StringAttribute{
					Optional:    true,
					Description: "IP address template to attach to logs",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
				"mac": schema.StringAttribute{
					Optional:    true,
					Description: "MAC address template to attach to logs",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
			},
		},
		"log_construction_scheme": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("explicit"),
			Description: "How to construct the log message",
			Validators: []validator.String{
				stringvalidator.OneOf(MapKeys(log_construction_schemes)...),
			},
		},
		"explicit_scheme_options": schema.SingleNestedAttribute{
			Optional:    true,
			Description: "Log construction options for the explicit scheme",
			Attributes: map[string]schema.Attribute{
				"line": schema.StringAttribute{
					Optional: true,
					Description: "Template or field reference to use as the log line. " +
						"Field reference can point to an object.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(1024),
					},
				},
				"meta_field": schema.StringAttribute{
					Optional:    true,
					Description: "Field containing the meta object for the log",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
				"app": schema.StringAttribute{
					Optional:    true,
					Description: "App name template to attach to logs",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
				"file": schema.StringAttribute{
					Optional:    true,
					Description: "File name template to attach to logs",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
				"timestamp_field": schema.StringAttribute{
					Optional:    true,
					Description: "Field containing the timestamp for the log, for example ._ts",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
				"env": schema.StringAttribute{
					Optional:    true,
					Description: "Environment name template to attach to logs",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.LengthAtMost(512),
					},
				},
			},
		},
	}, nil),
}

func MezmoDestinationFromModel(plan *MezmoDestinationModel, previousState *MezmoDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "mezmo",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":   plan.AckEnabled.ValueBool(),
				"mezmo_host":    plan.Host.ValueString(),
				"ingestion_key": plan.IngestionKey.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	if !plan.Query.IsNull() {
		plan_map := plan.Query.Attributes()
		query_options := make(map[string]any)
		SetOptionalStringFromAttributeMap(query_options, plan_map, "hostname", "ip", "mac")

		if plan_map["tags"] != nil && !plan_map["tags"].IsNull() {
			query_options["tags"] = StringListValueToStringSlice(GetAttributeValue[List](plan_map, "tags"))
		}

		component.UserConfig["query"] = query_options
	}

	log_scheme := plan.LogConstructionScheme.ValueString()
	message_options := map[string]any{
		"scheme": log_construction_schemes[log_scheme],
	}

	if !plan.ExplicitSchemeOptions.IsNull() {
		if log_scheme == "explicit" {
			plan_map := plan.ExplicitSchemeOptions.Attributes()
			SetOptionalStringFromAttributeMap(message_options, plan_map,
				"line", "app", "file", "env", "meta_field", "timestamp_field")
		} else {
			dd.AddError(
				"Invalid log constructions options",
				"Plan cannot define explicit scheme options when using message pass-through")
		}
	}

	component.UserConfig["message"] = message_options

	return &component, dd
}

func MezmoDestinationToModel(plan *MezmoDestinationModel, component *Destination) {
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
	plan.Host = StringValue(component.UserConfig["mezmo_host"].(string))
	plan.IngestionKey = StringValue(component.UserConfig["ingestion_key"].(string))

	if component.UserConfig["query"] != nil {
		component_map, _ := component.UserConfig["query"].(map[string]any)
		if len(component_map) > 0 {
			plan_map := map[string]attr.Value{
				"hostname": StringValue(component_map["hostname"].(string)),
				"tags":     ListNull(StringType),
				"ip":       StringNull(),
				"mac":      StringNull(),
			}
			if component_map["tags"] != nil {
				plan_map["tags"] = SliceToStringListValue(component_map["tags"].([]any))
			}
			SetOptionalAttributeStringFromMap(plan_map, component_map, "ip", "mac")

			attrTypes := plan.Query.AttributeTypes(context.Background())
			if len(attrTypes) == 0 {
				attrTypes = MezmoDestinationResourceSchema.Attributes["query"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.Query = basetypes.NewObjectValueMust(attrTypes, plan_map)
		}
	}

	if component.UserConfig["message"] != nil {
		component_map, _ := component.UserConfig["message"].(map[string]any)
		plan.LogConstructionScheme = StringValue(FindKey(log_construction_schemes, component_map["scheme"].(string)))
		plan_map := map[string]attr.Value{
			"line":            StringNull(),
			"app":             StringNull(),
			"file":            StringNull(),
			"env":             StringNull(),
			"meta_field":      StringNull(),
			"timestamp_field": StringNull(),
		}
		SetOptionalAttributeStringFromMap(plan_map, component_map,
			"line", "app", "file", "env", "meta_field", "timestamp_field")
		has_explicit_options := false
		for _, v := range plan_map {
			if !v.IsNull() {
				has_explicit_options = true
				break
			}
		}
		if has_explicit_options {
			attrTypes := plan.ExplicitSchemeOptions.AttributeTypes(context.Background())
			if len(attrTypes) == 0 {
				attrTypes = MezmoDestinationResourceSchema.Attributes["explicit_scheme_options"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.ExplicitSchemeOptions = basetypes.NewObjectValueMust(
				attrTypes,
				plan_map)
		}
	}
}
