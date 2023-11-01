package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const SQS_SOURCE_TYPE_NAME = "sqs"
const SQS_SOURCE_NODE_NAME = SQS_SOURCE_TYPE_NAME

type SQSSourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline_id"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	GenerationId Int64  `tfsdk:"generation_id"`
	QueueUrl     String `tfsdk:"queue_url" user_config:"true"`
	Auth         Object `tfsdk:"auth" user_config:"true"`
	Region       String `tfsdk:"region" user_config:"true"`
}

var SQSSourceResourceSchema = schema.Schema{
	Description: "Collect messages from AWS SQS",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"queue_url": schema.StringAttribute{
			Required:    true,
			Description: "The URL of an AWS SQS queue",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(7), // http://
				stringvalidator.LengthAtMost(128),
			},
		},
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
			Description: "The name of the source's AWS region",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
	}, nil),
}

func SQSSourceFromModel(plan *SQSSourceModel, previousState *SQSSourceModel) (*Source, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	auth := plan.Auth.Attributes()
	auth_access_key_id, _ := auth["access_key_id"].(basetypes.StringValue)
	auth_secret_access_key, _ := auth["secret_access_key"].(basetypes.StringValue)
	component := Source{
		BaseNode: BaseNode{
			Type:        SQS_SOURCE_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig: map[string]any{
				"region":    plan.Region.ValueString(),
				"queue_url": plan.QueueUrl.ValueString(),
				"auth": map[string]string{
					"access_key_id":     auth_access_key_id.ValueString(),
					"secret_access_key": auth_secret_access_key.ValueString(),
				},
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func SQSSourceToModel(plan *SQSSourceModel, component *Source) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}

	// Required properties will always be present in the API response
	region, _ := component.UserConfig["region"].(string)
	plan.Region = StringValue(region)

	queueUrl, _ := component.UserConfig["queue_url"].(string)
	plan.QueueUrl = StringValue(queueUrl)

	if component.UserConfig["auth"] != nil {
		values, _ := component.UserConfig["auth"].(map[string]any)
		if len(values) > 0 {
			objT := plan.Auth.AttributeTypes(context.Background())
			if len(objT) == 0 {
				objT = SQSSourceResourceSchema.Attributes["auth"].GetType().(basetypes.ObjectType).AttrTypes
			}
			plan.Auth = basetypes.NewObjectValueMust(objT, modelutils.MapAnyToMapValues(values))
		}
	}

	plan.GenerationId = Int64Value(component.GenerationId)
}
