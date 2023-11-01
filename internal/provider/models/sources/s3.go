package sources

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
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const S3_SOURCE_TYPE_NAME = "s3"
const S3_SOURCE_NODE_NAME = S3_SOURCE_TYPE_NAME

type S3SourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	GenerationId Int64  `tfsdk:"generation_id"`
	Auth         Object `tfsdk:"auth" user_config:"true"`
	Region       String `tfsdk:"region" user_config:"true"`
	SqsQueueUrl  String `tfsdk:"sqs_queue_url" user_config:"true"`
	Compression  String `tfsdk:"compression" user_config:"true"`
}

var S3SourceResourceSchema = schema.Schema{
	Description: "Represents an S3 pull source.",
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
		"sqs_queue_url": schema.StringAttribute{
			Required:    true,
			Description: "The URL of a AWS SQS queue configured to receive S3 bucket notifications",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(7)},
		},
		"compression": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("auto"),
			Description: "The compression format of the S3 objects",
			Validators:  []validator.String{stringvalidator.OneOf([]string{"auto", "gzip", "none", "zstd"}...)},
		},
	}, nil),
}

func S3SourceFromModel(plan *S3SourceModel, previousState *S3SourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	auth := plan.Auth.Attributes()
	component := Source{
		BaseNode: BaseNode{
			Type:        S3_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"region":        plan.Region.ValueString(),
				"sqs_queue_url": plan.SqsQueueUrl.ValueString(),
				"auth": map[string]string{
					"access_key_id":     GetAttributeValue[String](auth, "access_key_id").ValueString(),
					"secret_access_key": GetAttributeValue[String](auth, "secret_access_key").ValueString(),
				},
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

func S3SourceToModel(plan *S3SourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	if component.UserConfig["region"] != nil {
		value, _ := component.UserConfig["region"].(string)
		plan.Region = StringValue(value)
	}
	if component.UserConfig["sqs_queue_url"] != nil {
		value, _ := component.UserConfig["sqs_queue_url"].(string)
		plan.SqsQueueUrl = StringValue(value)
	}
	if component.UserConfig["compression"] != nil {
		value, _ := component.UserConfig["compression"].(string)
		plan.Compression = StringValue(value)
	}
	if component.UserConfig["auth"] != nil {
		values, _ := component.UserConfig["auth"].(map[string]any)
		if len(values) > 0 {
			objT := plan.Auth.AttributeTypes(context.Background())
			if len(objT) == 0 {
				objT = S3SourceResourceSchema.Attributes["auth"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.Auth = basetypes.NewObjectValueMust(objT, modelutils.MapAnyToMapValues(values))
		}
	}

	plan.GenerationId = Int64Value(component.GenerationId)
}
