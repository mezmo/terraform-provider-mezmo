package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/processors"
)

func NewAggregateProcessorResource() resource.Resource {
	return &ProcessorResource[AggregateProcessorModel]{
		typeName:          AGGREGATE_PROCESSOR_TYPE_NAME,
		nodeName:          AGGREGATE_PROCESSOR_NODE_NAME,
		fromModelFunc:     AggregateProcessorFromModel,
		toModelFunc:       AggregateProcessorToModel,
		getIdFunc:         func(m *AggregateProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AggregateProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            AggregateProcessorResourceSchema,
	}
}

func NewDedupeProcessorResource() resource.Resource {
	return &ProcessorResource[DedupeProcessorModel]{
		typeName:          DEDUPE_PROCESSOR_TYPE_NAME,
		nodeName:          DEDUPE_PROCESSOR_NODE_NAME,
		fromModelFunc:     DedupeProcessorFromModel,
		toModelFunc:       DedupeProcessorToModel,
		getIdFunc:         func(m *DedupeProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DedupeProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DedupeProcessorResourceSchema,
	}
}

func NewDropFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[DropFieldsProcessorModel]{
		typeName:          DROP_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          DROP_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     DropFieldsProcessorFromModel,
		toModelFunc:       DropFieldsProcessorToModel,
		getIdFunc:         func(m *DropFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DropFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DropFieldsProcessorResourceSchema,
	}
}

func NewFlattenFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[FlattenFieldsProcessorModel]{
		typeName:          FLATTEN_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          FLATTEN_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     FlattenFieldsProcessorFromModel,
		toModelFunc:       FlattenFieldsProcessorToModel,
		getIdFunc:         func(m *FlattenFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FlattenFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            FlattenFieldsProcessorResourceSchema,
	}
}

func NewMapFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[MapFieldsProcessorModel]{
		typeName:          MAP_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          MAP_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     MapFieldsProcessorFromModel,
		toModelFunc:       MapFieldsProcessorToModel,
		getIdFunc:         func(m *MapFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MapFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            MapFieldsProcessorResourceSchema,
	}
}

func NewSampleProcessorResource() resource.Resource {
	return &ProcessorResource[SampleProcessorModel]{
		typeName:          SAMPLE_PROCESSOR_TYPE_NAME,
		nodeName:          SAMPLE_PROCESSOR_NODE_NAME,
		fromModelFunc:     SampleProcessorFromModel,
		toModelFunc:       SampleProcessorToModel,
		getIdFunc:         func(m *SampleProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SampleProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            SampleProcessorResourceSchema,
	}
}

func NewStringifyProcessorResource() resource.Resource {
	return &ProcessorResource[StringifyProcessorModel]{
		typeName:          STRINGIFY_PROCESSOR_TYPE_NAME,
		nodeName:          STRINGIFY_PROCESSOR_NODE_NAME,
		fromModelFunc:     StringifyProcessorFromModel,
		toModelFunc:       StringifyProcessorToModel,
		getIdFunc:         func(m *StringifyProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *StringifyProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            StringifyProcessorResourceSchema,
	}
}

func NewScriptExecutionProcessorResource() resource.Resource {
	return &ProcessorResource[ScriptExecutionProcessorModel]{
		typeName:          SCRIPT_EXECUTION_PROCESSOR_TYPE_NAME,
		nodeName:          SCRIPT_EXECUTION_PROCESSOR_NODE_NAME,
		fromModelFunc:     ScriptExecutionProcessorFromModel,
		toModelFunc:       ScriptExecutionProcessorToModel,
		getIdFunc:         func(m *ScriptExecutionProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ScriptExecutionProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ScriptExecutionProcessorResourceSchema,
	}
}

func NewUnrollProcessorResource() resource.Resource {
	return &ProcessorResource[UnrollProcessorModel]{
		typeName:          UNROLL_PROCESSOR_TYPE_NAME,
		nodeName:          UNROLL_PROCESSOR_NODE_NAME,
		fromModelFunc:     UnrollProcessorFromModel,
		toModelFunc:       UnrollProcessorToModel,
		getIdFunc:         func(m *UnrollProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *UnrollProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            UnrollProcessorResourceSchema,
	}
}

func NewCompactFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[CompactFieldsProcessorModel]{
		typeName:          COMPACT_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          COMPACT_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     CompactFieldsProcessorFromModel,
		toModelFunc:       CompactFieldsProcessorToModel,
		getIdFunc:         func(m *CompactFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *CompactFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            CompactFieldsProcessorResourceSchema,
	}
}

func NewDecryptFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[DecryptFieldsProcessorModel]{
		typeName:          DECRYPT_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          DECRYPT_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     DecryptFieldsProcessorFromModel,
		toModelFunc:       DecryptFieldsProcessorToModel,
		getIdFunc:         func(m *DecryptFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DecryptFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DecryptFieldsProcessorResourceSchema,
	}
}

func NewEncryptFieldsProcessorResource() resource.Resource {
	return &ProcessorResource[EncryptFieldsProcessorModel]{
		typeName:          ENCRYPT_FIELDS_PROCESSOR_TYPE_NAME,
		nodeName:          ENCRYPT_FIELDS_PROCESSOR_NODE_NAME,
		fromModelFunc:     EncryptFieldsProcessorFromModel,
		toModelFunc:       EncryptFieldsProcessorToModel,
		getIdFunc:         func(m *EncryptFieldsProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *EncryptFieldsProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            EncryptFieldsProcessorResourceSchema,
	}
}

func NewParseProcessorResource() resource.Resource {
	return &ProcessorResource[ParseProcessorModel]{
		typeName:          PARSE_PROCESSOR_TYPE_NAME,
		nodeName:          PARSE_PROCESSOR_NODE_NAME,
		fromModelFunc:     ParseProcessorFromModel,
		toModelFunc:       ParseProcessorToModel,
		getIdFunc:         func(m *ParseProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ParseProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ParseProcessorResourceSchema,
	}
}

func NewReduceProcessorResource() resource.Resource {
	return &ProcessorResource[ReduceProcessorModel]{
		typeName:          REDUCE_PROCESSOR_TYPE_NAME,
		nodeName:          REDUCE_PROCESSOR_NODE_NAME,
		fromModelFunc:     ReduceProcessorFromModel,
		toModelFunc:       ReduceProcessorToModel,
		getIdFunc:         func(m *ReduceProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ReduceProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ReduceProcessorResourceSchema,
	}
}

func NewRouteProcessorResource() resource.Resource {
	return &ProcessorResource[RouteProcessorModel]{
		typeName:          ROUTE_PROCESSOR_TYPE_NAME,
		nodeName:          ROUTE_PROCESSOR_NODE_NAME,
		fromModelFunc:     RouteProcessorFromModel,
		toModelFunc:       RouteProcessorToModel,
		getIdFunc:         func(m *RouteProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *RouteProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            RouteProcessorResourceSchema,
	}
}

func NewParseSequentiallyProcessorResource() resource.Resource {
	return &ProcessorResource[ParseSequentiallyProcessorModel]{
		typeName:          PARSE_SEQUENTIALLY_PROCESSOR_TYPE_NAME,
		nodeName:          PARSE_SEQUENTIALLY_PROCESSOR_NODE_NAME,
		fromModelFunc:     ParseSequentiallyProcessorFromModel,
		toModelFunc:       ParseSequentiallyProcessorToModel,
		getIdFunc:         func(m *ParseSequentiallyProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ParseSequentiallyProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ParseSequentiallyProcessorResourceSchema,
	}
}

func NewMetricsTagCardinalityLimitProcessorResource() resource.Resource {
	return &ProcessorResource[MetricsTagCardinalityLimitProcessorModel]{
		typeName:          METRICS_TAG_CARDINALITY_LIMIT_PROCESSOR_TYPE_NAME,
		nodeName:          METRICS_TAG_LIMIT_PROCESSOR_NODE_NAME,
		fromModelFunc:     MetricsTagCardinalityLimitProcessorFromModel,
		toModelFunc:       MetricsTagCardinalityLimitProcessorToModel,
		getIdFunc:         func(m *MetricsTagCardinalityLimitProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MetricsTagCardinalityLimitProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            MetricsTagCardinalityLimitProcessorResourceSchema,
	}
}

func NewEventToMetricProcessorResource() resource.Resource {
	return &ProcessorResource[EventToMetricProcessorModel]{
		typeName:          EVENT_TO_METRIC_PROCESSOR_TYPE_NAME,
		nodeName:          EVENT_TO_METRIC_PROCESSOR_NODE_NAME,
		fromModelFunc:     EventToMetricProcessorFromModel,
		toModelFunc:       EventToMetricProcessorToModel,
		getIdFunc:         func(m *EventToMetricProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *EventToMetricProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            EventToMetricProcessorResourceSchema,
	}
}

func NewFilterProcessorResource() resource.Resource {
	return &ProcessorResource[FilterProcessorModel]{
		typeName:          FILTER_PROCESSOR_TYPE_NAME,
		nodeName:          FILTER_PROCESSOR_NODE_NAME,
		fromModelFunc:     FilterProcessorFromModel,
		toModelFunc:       FilterProcessorToModel,
		getIdFunc:         func(m *FilterProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FilterProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            FilterProcessorResourceSchema,
	}
}

func NewSetTimestampProcessorResource() resource.Resource {
	return &ProcessorResource[SetTimestampProcessorModel]{
		typeName:          SET_TIMESTAMP_PROCESSOR_TYPE_NAME,
		nodeName:          SET_TIMESTAMP_PROCESSOR_NODE_NAME,
		fromModelFunc:     SetTimestampProcessorFromModel,
		toModelFunc:       SetTimestampProcessorToModel,
		getIdFunc:         func(m *SetTimestampProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SetTimestampProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            SetTimestampProcessorResourceSchema,
	}
}

func NewThrottleProcessorResource() resource.Resource {
	return &ProcessorResource[ThrottleProcessorModel]{
		typeName:          THROTTLE_PROCESSOR_TYPE_NAME,
		nodeName:          THROTTLE_PROCESSOR_NODE_NAME,
		fromModelFunc:     ThrottleProcessorFromModel,
		toModelFunc:       ThrottleProcessorToModel,
		getIdFunc:         func(m *ThrottleProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ThrottleProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            ThrottleProcessorResourceSchema,
	}
}

func NewDataProfilerProcessorResource() resource.Resource {
	return &ProcessorResource[DataProfilerProcessorModel]{
		typeName:          DATA_PROFILER_PROCESSOR_TYPE_NAME,
		nodeName:          DATA_PROFILER_PROCESSOR_NODE_NAME,
		fromModelFunc:     DataProfilerProcessorFromModel,
		toModelFunc:       DataProfilerProcessorToModel,
		getIdFunc:         func(m *DataProfilerProcessorModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DataProfilerProcessorModel) basetypes.StringValue { return m.PipelineId },
		schema:            DataProfilerProcessorResourceSchema,
	}
}
