package processors

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

const DATA_PROFILER_PROCESSOR_TYPE_NAME = "data_profiler"
const DATA_PROFILER_PROCESSOR_NODE_NAME = "data-profiler"

type DataProfilerProcessorModel struct {
	Id           StringValue `tfsdk:"id"`
	PipelineId   StringValue `tfsdk:"pipeline_id"`
	Title        StringValue `tfsdk:"title"`
	Description  StringValue `tfsdk:"description"`
	Inputs       ListValue   `tfsdk:"inputs"`
	GenerationId Int64Value  `tfsdk:"generation_id"`
	AppFields    ListValue   `tfsdk:"app_fields" user_config:"true"`
	HostFields   ListValue   `tfsdk:"host_fields" user_config:"true"`
	LevelFields  ListValue   `tfsdk:"level_fields" user_config:"true"`
	LineFields   ListValue   `tfsdk:"line_fields" user_config:"true"`
	LabelFields  ListValue   `tfsdk:"label_fields" user_config:"true"`
}

var DataProfilerProcessorResourceSchema = schema.Schema{
	Description: "Profile the data sent through your pipeline, generating annotations and saving profile data.",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"app_fields":   ListAttributeSchemaGenerator("app"),
		"host_fields":  ListAttributeSchemaGenerator("host"),
		"level_fields": ListAttributeSchemaGenerator("level"),
		"line_fields":  ListAttributeSchemaGenerator("line"),
		"label_fields": schema.ListAttribute{
			ElementType: basetypes.StringType{},
			Optional:    true,
			Description: "A list of paths to look for the label value.",
			Validators: []validator.List{
				listvalidator.SizeAtMost(10),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				listvalidator.ValueStringsAre(stringvalidator.LengthAtMost(512)),
			},
		},
	}),
}

func ListAttributeSchemaGenerator(attributeName string) schema.ListAttribute {
	return schema.ListAttribute{
		ElementType: basetypes.StringType{},
		Required:    true,
		Description: fmt.Sprintf("A list of paths to look for the %s value.", attributeName),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.SizeAtMost(10),
			listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			listvalidator.ValueStringsAre(stringvalidator.LengthAtMost(512)),
		},
	}
}

func DataProfilerProcessorFromModel(plan *DataProfilerProcessorModel, previousState *DataProfilerProcessorModel) (*Processor, diag.Diagnostics) {
	dd := diag.Diagnostics{}
	component := Processor{
		BaseNode: BaseNode{
			Type:        DATA_PROFILER_PROCESSOR_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			UserConfig:  map[string]any{},
		},
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}
	component.Inputs = StringListValueToStringSlice(plan.Inputs)

	component.UserConfig["app_fields"] = StringListValueToStringSlice(plan.AppFields)
	component.UserConfig["host_fields"] = StringListValueToStringSlice(plan.HostFields)
	component.UserConfig["level_fields"] = StringListValueToStringSlice(plan.LevelFields)
	component.UserConfig["line_fields"] = StringListValueToStringSlice(plan.LineFields)
	if !plan.LabelFields.IsNull() {
		component.UserConfig["label_fields"] = StringListValueToStringSlice(plan.LabelFields)
	}

	return &component, dd
}

func DataProfilerProcessorToModel(plan *DataProfilerProcessorModel, component *Processor) {
	plan.Id = NewStringValue(component.Id)
	if component.Title != "" {
		plan.Title = NewStringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = NewStringValue(component.Description)
	}
	plan.GenerationId = NewInt64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)

	plan.AppFields = SliceToStringListValue(component.UserConfig["app_fields"].([]any))
	plan.HostFields = SliceToStringListValue(component.UserConfig["host_fields"].([]any))
	plan.LevelFields = SliceToStringListValue(component.UserConfig["level_fields"].([]any))
	plan.LineFields = SliceToStringListValue(component.UserConfig["line_fields"].([]any))
	if component.UserConfig["label_fields"] != nil {
		plan.LabelFields = SliceToStringListValue(component.UserConfig["label_fields"].([]any))
	}

}
