package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/processors"
)

func NewDedupeProcessorResource() resource.Resource {
	return &ProcessorResource[DedupeProcessorModel]{
		typeName:          "dedupe",
		fromModelFunc:     DedupeProcessorFromModel,
		toModelFunc:       DedupeProcessorToModel,
		getIdFunc:         func(m *DedupeProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DedupeProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DedupeProcessorResourceSchema,
	}
}

func NewDropFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[DropFieldsProcessorModel]{
		typeName:          "drop_fields",
		fromModelFunc:     DropFieldsProcessorFromModel,
		toModelFunc:       DropFieldsProcessorToModel,
		getIdFunc:         func(m *DropFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DropFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DropFieldsProcessorResourceSchema,
	}
}

func NewFlattenFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[FlattenFieldsProcessorModel]{
		typeName:          "flatten_fields",
		fromModelFunc:     FlattenFieldsProcessorFromModel,
		toModelFunc:       FlattenFieldsProcessorToModel,
		getIdFunc:         func(m *FlattenFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FlattenFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            FlattenFieldsProcessorResourceSchema,
	}
}

func NewMapFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[MapFieldsProcessorModel]{
		typeName:          "map_fields",
		fromModelFunc:     MapFieldsProcessorFromModel,
		toModelFunc:       MapFieldsProcessorToModel,
		getIdFunc:         func(m *MapFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MapFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            MapFieldsProcessorResourceSchema,
	}
}

func NewSampleProcessorResource() resource.Resource {
	return &ProcessorResource[SampleProcessorModel]{
		typeName:          "sample",
		fromModelFunc:     SampleProcessorFromModel,
		toModelFunc:       SampleProcessorToModel,
		getIdFunc:         func(m *SampleProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SampleProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            SampleProcessorResourceSchema,
	}
}

func NewStringifyProcessorResource() resource.Resource {
	return &ProcessorResource[StringifyProcessorModel]{
		typeName:          "stringify",
		fromModelFunc:     StringifyProcessorFromModel,
		toModelFunc:       StringifyProcessorToModel,
		getIdFunc:         func(m *StringifyProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *StringifyProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            StringifyProcessorResourceSchema,
	}
}

func NewScriptExecutionProcessorResource() resource.Resource {
	return &ProcessorResource[ScriptExecutionProcessorModel]{
		typeName:          "script_execution",
		fromModelFunc:     ScriptExecutionProcessorFromModel,
		toModelFunc:       ScriptExecutionProcessorToModel,
		getIdFunc:         func(m *ScriptExecutionProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ScriptExecutionProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ScriptExecutionProcessorResourceSchema,
	}
}

func NewUnrollProcessorResource() resource.Resource {
	return &ProcessorResource[UnrollProcessorModel]{
		typeName:          "unroll",
		fromModelFunc:     UnrollProcessorFromModel,
		toModelFunc:       UnrollProcessorToModel,
		getIdFunc:         func(m *UnrollProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *UnrollProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            UnrollProcessorResourceSchema,
	}
}

func NewCompactFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[CompactFieldsProcessorModel]{
		typeName:          "compact_fields",
		fromModelFunc:     CompactFieldsProcessorFromModel,
		toModelFunc:       CompactFieldsProcessorToModel,
		getIdFunc:         func(m *CompactFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *CompactFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            CompactFieldsProcessorResourceSchema,
	}
}

func NewDecryptFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[DecryptFieldsProcessorModel]{
		typeName:          "decrypt_fields",
		fromModelFunc:     DecryptFieldsProcessorFromModel,
		toModelFunc:       DecryptFieldsProcessorToModel,
		getIdFunc:         func(m *DecryptFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DecryptFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DecryptFieldsProcessorResourceSchema,
	}
}

func NewEncryptFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[EncryptFieldsProcessorModel]{
		typeName:          "encrypt_fields",
		fromModelFunc:     EncryptFieldsProcessorFromModel,
		toModelFunc:       EncryptFieldsProcessorToModel,
		getIdFunc:         func(m *EncryptFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *EncryptFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            EncryptFieldsProcessorResourceSchema,
	}
}

func NewParseProcessorResource() resource.Resource {
	return &ProcessorResource[ParseProcessorModel]{
		typeName:          "parse",
		fromModelFunc:     ParseProcessorFromModel,
		toModelFunc:       ParseProcessorToModel,
		getIdFunc:         func(m *ParseProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ParseProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ParseProcessorResourceSchema,
	}
}

func NewReduceProcessorResource() resource.Resource {
	return &ProcessorResource[ReduceProcessorModel]{
		typeName:          "reduce",
		fromModelFunc:     ReduceProcessorFromModel,
		toModelFunc:       ReduceProcessorToModel,
		getIdFunc:         func(m *ReduceProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ReduceProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ReduceProcessorResourceSchema,
	}
}

func NewRouteProcessorResource() resource.Resource {
	return &ProcessorResource[RouteProcessorModel]{
		typeName:          RouteProcessorName,
		fromModelFunc:     RouteProcessorFromModel,
		toModelFunc:       RouteProcessorToModel,
		getIdFunc:         func(m *RouteProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *RouteProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            RouteProcessorResourceSchema,
	}
}

func NewParseSequentiallyProcessorResource() resource.Resource {
	return &ProcessorResource[ParseSequentiallyProcessorModel]{
		typeName:          ParseSequentiallyProcessorName,
		fromModelFunc:     ParseSequentiallyProcessorFromModel,
		toModelFunc:       ParseSequentiallyProcessorToModel,
		getIdFunc:         func(m *ParseSequentiallyProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ParseSequentiallyProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ParseSequentiallyProcessorResourceSchema,
	}
}
