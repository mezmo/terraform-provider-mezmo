package destinations

import (
	"context"

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

type GcpCloudStorageDestinationModel struct {
	Id                  String `tfsdk:"id"`
	PipelineId          String `tfsdk:"pipeline_id"`
	Title               String `tfsdk:"title"`
	Description         String `tfsdk:"description"`
	Inputs              List   `tfsdk:"inputs"`
	GenerationId        Int64  `tfsdk:"generation_id"`
	Encoding            String `tfsdk:"encoding" user_config:"true"`
	Bucket              String `tfsdk:"bucket" user_config:"true"`
	Compression         String `tfsdk:"compression" user_config:"true"`
	BucketPrefix        String `tfsdk:"bucket_prefix" user_config:"true"`
	Auth                Object `tfsdk:"auth" user_config:"true"`
	AckEnabled          Bool   `tfsdk:"ack_enabled" user_config:"true"`
	BatchTimeoutSeconds Int64  `tfsdk:"batch_timeout_secs" user_config:"true"`
}

var GcpCloudStorageResourceSchema = schema.Schema{
	Description: "Publish log events to GCP Cloud Storage",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"encoding": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Dictates how the data will be serialized before storing.",
			Default:     stringdefault.StaticString("text"),
			Validators: []validator.String{
				stringvalidator.OneOf("json", "text"),
			},
		},
		"bucket": schema.StringAttribute{
			Required:    true,
			Description: "The name of the bucket in GCP where the data will be stored.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"compression": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The compression strategy used on the encoded data prior to sending.",
			Default:     stringdefault.StaticString("none"),
			Validators: []validator.String{
				stringvalidator.OneOf("gzip", "none"),
			},
		},
		"bucket_prefix": schema.StringAttribute{
			Optional:    true,
			Computed:    false,
			Description: "The prefix applied to the bucket name, giving the appearance of having directories.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"auth": schema.SingleNestedAttribute{
			Required:    true,
			Description: "Configure GCP Cloud Storage authentication",
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:    true,
					Description: "The type of authentication to use.",
					Validators: []validator.String{
						stringvalidator.OneOf("api_key", "credentials_json"),
					},
				},
				"value": schema.StringAttribute{
					Required:    true,
					Sensitive:   true,
					Description: "Authentication secret value.",
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
			},
		},
	}, []string{"batch_timeout_secs"}),
}

func GcpCloudStorageDestinationFromModel(plan *GcpCloudStorageDestinationModel, previousState *GcpCloudStorageDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        "gcp-cloud-storage",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":        plan.AckEnabled.ValueBool(),
				"batch_timeout_secs": plan.BatchTimeoutSeconds.ValueInt64(),
				"encoding":           plan.Encoding.ValueString(),
				"bucket":             plan.Bucket.ValueString(),
				"compression":        plan.Compression.ValueString(),
				"bucket_prefix":      plan.BucketPrefix.ValueString(),
			},
		},
	}

	auth := plan.Auth.Attributes()
	component.UserConfig["auth"] = GetAttributeValue[String](auth, "type").ValueString()
	if component.UserConfig["auth"] == "api_key" {
		component.UserConfig["api_key"] = GetAttributeValue[String](auth, "value").ValueString()
	} else {
		component.UserConfig["credentials_json"] = GetAttributeValue[String](auth, "value").ValueString()
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func GcpCloudStorageDestinationToModel(plan *GcpCloudStorageDestinationModel, component *Destination) {
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
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.BucketPrefix = StringValue(component.UserConfig["bucket_prefix"].(string))
	plan.Bucket = StringValue(component.UserConfig["bucket"].(string))
	plan.BatchTimeoutSeconds = Int64Value(int64(component.UserConfig["batch_timeout_secs"].(float64)))

	authType, _ := component.UserConfig["auth"].(string)
	plan.Auth = basetypes.NewObjectValueMust(
		plan.Auth.AttributeTypes(context.Background()),
		map[string]attr.Value{
			"type":  StringValue(authType),
			"value": StringValue(component.UserConfig[authType].(string)),
		})
}
