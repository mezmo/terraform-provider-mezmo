package sources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

type S3SourceModel struct {
	Id           String `tfsdk:"id"`
	PipelineId   String `tfsdk:"pipeline"`
	Title        String `tfsdk:"title"`
	Description  String `tfsdk:"description"`
	Auth         Object `tfsdk:"auth"`
	Region       String `tfsdk:"region"`
	SqsQueueUrl  String `tfsdk:"sqs_queue_url"`
	GenerationId Int64  `tfsdk:"generation_id"`
}

func S3SourceResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pipeline": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"title": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(256),
				},
			},
			"description": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"auth": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"access_key_id": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"secret_access_key": schema.StringAttribute{
						Required:   true,
						Sensitive:  true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
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
			"generation_id": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func S3SourceFromModel(model *S3SourceModel, previousState *S3SourceModel) *Component {
	auth := model.Auth.Attributes()
	auth_access_key_id, _ := auth["access_key_id"].(basetypes.StringValue)
	auth_secret_access_key, _ := auth["secret_access_key"].(basetypes.StringValue)
	component := Component{
		Type:        "s3",
		Title:       model.Title.ValueString(),
		Description: model.Description.ValueString(),
		UserConfig: map[string]any{
			"region":        model.Region.ValueString(),
			"sqs_queue_url": model.SqsQueueUrl.ValueString(),
			"auth": map[string]string{
				"access_key_id":     auth_access_key_id.ValueString(),
				"secret_access_key": auth_secret_access_key.ValueString(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component
}

func S3SourceToModel(model *S3SourceModel, component *Component) {
	model.Id = StringValue(component.Id)
	if component.Title != "" {
		model.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		model.Description = StringValue(component.Description)
	}
	if component.UserConfig["region"] != nil {
		value, _ := component.UserConfig["region"].(string)
		model.Region = StringValue(value)
	}
	if component.UserConfig["sqs_queue_url"] != nil {
		value, _ := component.UserConfig["sqs_queue_url"].(string)
		model.SqsQueueUrl = StringValue(value)
	}
	if component.UserConfig["auth"] != nil {
		values, _ := component.UserConfig["auth"].(map[string]string)
		if len(values) > 0 {
			types := model.Auth.AttributeTypes(context.Background())
			model.Auth = basetypes.NewObjectValueMust(types, toAttributes(values))
		}
	}

	model.GenerationId = Int64Value(component.GenerationId)
}

func toAttributes(values map[string]string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(values))
	for k, v := range values {
		result[k] = StringValue(v)
	}
	return result
}
