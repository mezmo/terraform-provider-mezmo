package models

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
)

type AccessKeyResourceModel struct {
	Id       StringValue `tfsdk:"id"`
	Title    StringValue `tfsdk:"title" user_config:"true"`
	SourceId StringValue `tfsdk:"source_id" user_config:"true"`
	Key      StringValue `tfsdk:"key" user_config:"true"`
}

const AccessKeyType = "generated"

func AccessKeyResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the access key",
				Computed:    true,
			},
			"title": schema.StringAttribute{
				Description: "A descriptive title for the key and/or its use.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(512),
				},
			},
			"source_id": schema.StringAttribute{
				Description: "The uuid of the source (shared or not) for which the access key is created. ",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Description: "The cleartext key used for pipeline ingestion of the `source_id`." +
					"It is always a generated value meant for one-time consumption.",
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

// From terraform schema/model to a struct for sending to the API
func AccessKeyFromModel(plan *AccessKeyResourceModel) *AccessKey {
	accessKey := AccessKey{
		Title:          plan.Title.ValueString(),
		SharedSourceId: plan.SourceId.ValueString(),
		Type:           AccessKeyType,
	}
	if !plan.Id.IsUnknown() {
		accessKey.Id = plan.Id.ValueString()
	}
	// `Key` never needs to be sent to the api server. Leave it out.
	return &accessKey
}

// From an API response to a terraform model
func AccessKeyToModel(plan *AccessKeyResourceModel, accessKey *AccessKey) {
	plan.Id = NewStringValue(accessKey.Id)
	plan.Title = NewStringValue(accessKey.Title)
	plan.SourceId = NewStringValue(accessKey.SharedSourceId)
	plan.Key = NewStringValue(accessKey.Key)
}
