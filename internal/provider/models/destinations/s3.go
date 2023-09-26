package destinations

import (
	"context"

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

type S3DestinationModel struct {
	Id                  String `tfsdk:"id"`
	PipelineId          String `tfsdk:"pipeline_id"`
	Title               String `tfsdk:"title"`
	Description         String `tfsdk:"description"`
	Inputs              List   `tfsdk:"inputs"`
	GenerationId        Int64  `tfsdk:"generation_id"`
	AckEnabled          Bool   `tfsdk:"ack_enabled" user_config:"true"`
	BatchTimeoutSeconds Int64  `tfsdk:"batch_timeout_secs" user_config:"true"`
	Auth                Object `tfsdk:"auth" user_config:"true"`
	Region              String `tfsdk:"region" user_config:"true"`
	Bucket              String `tfsdk:"bucket" user_config:"true"`
	Prefix              String `tfsdk:"prefix" user_config:"true"`
	Encoding            String `tfsdk:"encoding" user_config:"true"`
	Compression         String `tfsdk:"compression" user_config:"true"`
}

func S3DestinationResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Publishes events as objects in AWS S3",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"auth": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Configures AWS authentication",
				Attributes: map[string]schema.Attribute{
					"access_key_id": schema.StringAttribute{
						Required:    true,
						Description: "The AWS access key id",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"secret_access_key": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The AWS secret access key",
						Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The name of the AWS region",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"bucket": schema.StringAttribute{
				Required:    true,
				Description: "The S3 bucket name. Do not include a leading s3:// or a trailing /",
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"prefix": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A prefix to apply to all object key names.",
				Default:     stringdefault.StaticString("/"),
				Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"encoding": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The encoding to apply to the data",
				Default:     stringdefault.StaticString("text"),
				Validators:  []validator.String{stringvalidator.OneOf("json", "ndjson", "text")},
			},
			"compression": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("none"),
				Description: "The compression format of the S3 objects",
				Validators:  []validator.String{stringvalidator.OneOf([]string{"gzip", "none"}...)},
			},
		}, []string{"batch_timeout_secs"}),
	}
}

func S3DestinationFromModel(plan *S3DestinationModel, previousState *S3DestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	auth := plan.Auth.Attributes()

	component := Destination{
		BaseNode: BaseNode{
			Type:        "s3",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":        plan.AckEnabled.ValueBool(),
				"batch_timeout_secs": plan.BatchTimeoutSeconds.ValueInt64(),
				"auth": map[string]string{
					"access_key_id":     GetAttributeValue[String](auth, "access_key_id").ValueString(),
					"secret_access_key": GetAttributeValue[String](auth, "secret_access_key").ValueString(),
				},
				"region":      plan.Region.ValueString(),
				"bucket":      plan.Bucket.ValueString(),
				"prefix":      plan.Prefix.ValueString(),
				"encoding":    plan.Encoding.ValueString(),
				"compression": plan.Compression.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func S3DestinationToModel(plan *S3DestinationModel, component *Destination) {
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
	plan.BatchTimeoutSeconds = Int64Value(int64(component.UserConfig["batch_timeout_secs"].(float64)))

	values, _ := component.UserConfig["auth"].(map[string]any)
	if len(values) > 0 {
		types := plan.Auth.AttributeTypes(context.Background())
		plan.Auth = basetypes.NewObjectValueMust(types, MapAnyToMapValues(values))
	}

	plan.Region = StringValue(component.UserConfig["region"].(string))
	plan.Bucket = StringValue(component.UserConfig["bucket"].(string))
	plan.Prefix = StringValue(component.UserConfig["prefix"].(string))
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.Compression = StringValue(component.UserConfig["compression"].(string))
}
