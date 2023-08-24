package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/transforms"
)

func NewDedupeTransformResource() resource.Resource {
	return &TransformResource[DedupeTransformModel]{
		typeName:          "dedupe",
		fromModelFunc:     DedupeTransformFromModel,
		toModelFunc:       DedupeTransformToModel,
		getIdFunc:         func(m *DedupeTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DedupeTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     DedupeTransformResourceSchema,
	}
}

func NewDropFieldsTransformResource() resource.Resource {
	return &TransformResource[DropFieldsTransformModel]{
		typeName:          "drop_fields",
		fromModelFunc:     DropFieldsTransformFromModel,
		toModelFunc:       DropFieldsTransformToModel,
		getIdFunc:         func(m *DropFieldsTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DropFieldsTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     DropFieldsTransformResourceSchema,
	}
}

func NewFlattenFieldsTransformResource() resource.Resource {
	return &TransformResource[FlattenFieldsTransformModel]{
		typeName:          "flatten_fields",
		fromModelFunc:     FlattenFieldsTransformFromModel,
		toModelFunc:       FlattenFieldsTransformToModel,
		getIdFunc:         func(m *FlattenFieldsTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FlattenFieldsTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     FlattenFieldsTransformResourceSchema,
	}
}

func NewSampleTransformResource() resource.Resource {
	return &TransformResource[SampleTransformModel]{
		typeName:          "sample",
		fromModelFunc:     SampleTransformFromModel,
		toModelFunc:       SampleTransformToModel,
		getIdFunc:         func(m *SampleTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SampleTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     SampleTransformResourceSchema,
	}
}

func NewStringifyTransformResource() resource.Resource {
	return &TransformResource[StringifyTransformModel]{
		typeName:          "stringify",
		fromModelFunc:     StringifyTransformFromModel,
		toModelFunc:       StringifyTransformToModel,
		getIdFunc:         func(m *StringifyTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *StringifyTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     StringifyTransformResourceSchema,
	}
}

func NewUnrollTransformResource() resource.Resource {
	return &TransformResource[UnrollTransformModel]{
		typeName:          "unroll",
		fromModelFunc:     UnrollTransformFromModel,
		toModelFunc:       UnrollTransformToModel,
		getIdFunc:         func(m *UnrollTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *UnrollTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     UnrollTransformResourceSchema,
	}
}

func NewCompactFieldsTransformResource() resource.Resource {
	return &TransformResource[CompactFieldsTransformModel]{
		typeName:          "compact_fields",
		fromModelFunc:     CompactFieldsTransformFromModel,
		toModelFunc:       CompactFieldsTransformToModel,
		getIdFunc:         func(m *CompactFieldsTransformModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *CompactFieldsTransformModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     CompactFieldsTransformResourceSchema,
	}
}
