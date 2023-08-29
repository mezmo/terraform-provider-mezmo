package transforms

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type EncryptFieldsTransformModel struct {
	Id             String `tfsdk:"id"`
	PipelineId     String `tfsdk:"pipeline_id"`
	Title          String `tfsdk:"title"`
	Description    String `tfsdk:"description"`
	Inputs         List   `tfsdk:"inputs"`
	GenerationId   Int64  `tfsdk:"generation_id"`
	Field          String `tfsdk:"field"`
	Algorithm      String `tfsdk:"algorithm"`
	Key            String `tfsdk:"key"`
	IvField        String `tfsdk:"iv_field"`
	EncodeRawBytes Bool   `tfsdk:"encode_raw_bytes"`
}

func EncryptFieldsTransformResourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Encrypts the value of the provided field",
		Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
			"field": schema.StringAttribute{
				Required:    true,
				Description: "Field to encrypt. The value of the field must be a primitive (string, number, boolean).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"algorithm": schema.StringAttribute{
				Required:    true,
				Description: "The encryption algorithm to use on the field",
				Validators: []validator.String{
					stringvalidator.OneOf(EncryptionAlgorithms...),
				},
			},
			"key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The encryption key",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(16),
					stringvalidator.LengthAtMost(32),
				},
			},
			"iv_field": schema.StringAttribute{
				Required: true,
				Description: "The field in which to store the generated initialization " +
					"vector, IV. Each encrypted value will have a unique IV.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"encode_raw_bytes": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				Description: "Encode the encrypted value and generated initialization " +
					"vector as Base64 text",
			},
		}),
	}
}

func EncryptFieldsTransformFromModel(plan *EncryptFieldsTransformModel, previousState *EncryptFieldsTransformModel) (*Transform, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Transform{
		BaseNode: BaseNode{
			Type:        "encrypt-fields",
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"field":            plan.Field.ValueString(),
				"algorithm":        plan.Algorithm.ValueString(),
				"key":              plan.Key.ValueString(),
				"iv_field":         plan.IvField.ValueString(),
				"encode_raw_bytes": plan.EncodeRawBytes.ValueBool(),
			},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func EncryptFieldsTransformToModel(plan *EncryptFieldsTransformModel, component *Transform) {
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
	plan.EncodeRawBytes = BoolValue(component.UserConfig["encode_raw_bytes"].(bool))
}