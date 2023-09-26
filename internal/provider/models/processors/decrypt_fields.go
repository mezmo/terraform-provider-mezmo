package processors

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type DecryptFieldsProcessorModel struct {
	Id             String `tfsdk:"id"`
	PipelineId     String `tfsdk:"pipeline_id"`
	Title          String `tfsdk:"title"`
	Description    String `tfsdk:"description"`
	Inputs         List   `tfsdk:"inputs"`
	GenerationId   Int64  `tfsdk:"generation_id"`
	Field          String `tfsdk:"field" user_config:"true"`
	Algorithm      String `tfsdk:"algorithm" user_config:"true"`
	Key            String `tfsdk:"key" user_config:"true"`
	IvField        String `tfsdk:"iv_field" user_config:"true"`
	DecodeRawBytes Bool   `tfsdk:"decode_raw_bytes" user_config:"true"`
}

func DecryptFieldsProcessorResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Decrypts the value of the provided field",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"field": schema.StringAttribute{
				Required:    true,
				Description: "Field to decrypt. The value of the field must be a string",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"algorithm": schema.StringAttribute{
				Required:    true,
				Description: "The algorithm with which the data was encrypted",
				Validators: []validator.String{
					stringvalidator.OneOf(EncryptionAlgorithms...),
				},
			},
			"key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The key/secret used to encrypt the value",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(16),
					stringvalidator.LengthAtMost(32),
				},
			},
			"iv_field": schema.StringAttribute{
				Required:    true,
				Description: "The field from which to read the initialization vector, IV",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"decode_raw_bytes": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "The field from which to read the initialization vector, IV",
			},
		}),
	}
}

func DecryptFieldsProcessorFromModel(plan *DecryptFieldsProcessorModel, previousState *DecryptFieldsProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        "decrypt-fields",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"field":            plan.Field.ValueString(),
				"algorithm":        plan.Algorithm.ValueString(),
				"key":              plan.Key.ValueString(),
				"iv_field":         plan.IvField.ValueString(),
				"decode_raw_bytes": plan.DecodeRawBytes.ValueBool(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func DecryptFieldsProcessorToModel(plan *DecryptFieldsProcessorModel, component *Processor) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)

	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.Field = StringValue(component.UserConfig["field"].(string))
	plan.Algorithm = StringValue(component.UserConfig["algorithm"].(string))
	plan.Key = StringValue(component.UserConfig["key"].(string))
	plan.IvField = StringValue(component.UserConfig["iv_field"].(string))
	plan.DecodeRawBytes = BoolValue(component.UserConfig["decode_raw_bytes"].(bool))
}
