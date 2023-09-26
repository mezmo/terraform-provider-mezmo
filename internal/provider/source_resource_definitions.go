package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/sources"
)

func NewDemoSourceResource() resource.Resource {
	return &SourceResource[DemoSourceModel]{
		typeName:          "demo",
		fromModelFunc:     DemoSourceFromModel,
		toModelFunc:       DemoSourceToModel,
		getIdFunc:         func(m *DemoSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DemoSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     DemoSourceResourceSchema,
	}
}

func NewAgentSourceResource() resource.Resource {
	return &SourceResource[AgentSourceModel]{
		typeName:          "agent",
		fromModelFunc:     AgentSourceFromModel,
		toModelFunc:       AgentSourceToModel,
		getIdFunc:         func(m *AgentSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AgentSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     AgentSourceResourceSchema,
	}
}

func NewKafkaSourceResource() resource.Resource {
	return &SourceResource[KafkaSourceModel]{
		typeName:          "kafka",
		fromModelFunc:     KafkaSourceFromModel,
		toModelFunc:       KafkaSourceToModel,
		getIdFunc:         func(m *KafkaSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KafkaSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     KafkaSourceResourceSchema,
	}
}

func NewPrometheusRemoteWriteSourceResource() resource.Resource {
	return &SourceResource[PrometheusRemoteWriteSourceModel]{
		typeName:          "prometheus_remote_write",
		fromModelFunc:     PrometheusRemoteWriteSourceFromModel,
		toModelFunc:       PrometheusRemoteWriteSourceToModel,
		getIdFunc:         func(m *PrometheusRemoteWriteSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *PrometheusRemoteWriteSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     PrometheusRemoteWriteSourceResourceSchema,
	}
}

func NewS3SourceResource() resource.Resource {
	return &SourceResource[S3SourceModel]{
		typeName:          "s3",
		fromModelFunc:     S3SourceFromModel,
		toModelFunc:       S3SourceToModel,
		getIdFunc:         func(m *S3SourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *S3SourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     S3SourceResourceSchema,
	}
}

func NewHttpSourceResource() resource.Resource {
	return &SourceResource[HttpSourceModel]{
		typeName:          "http",
		fromModelFunc:     HttpSourceFromModel,
		toModelFunc:       HttpSourceToModel,
		getIdFunc:         func(m *HttpSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HttpSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     HttpSourceResourceSchema,
	}
}

func NewSQSSourceResource() resource.Resource {
	return &SourceResource[SQSSourceModel]{
		typeName:          "sqs",
		fromModelFunc:     SQSSourceFromModel,
		toModelFunc:       SQSSourceToModel,
		getIdFunc:         func(m *SQSSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SQSSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     SQSSourceResourceSchema,
	}
}

func NewSplunkHecSourceResource() resource.Resource {
	return &SourceResource[SplunkHecSourceModel]{
		typeName:          "splunk_hec",
		fromModelFunc:     SplunkHecSourceFromModel,
		toModelFunc:       SplunkHecSourceToModel,
		getIdFunc:         func(m *SplunkHecSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SplunkHecSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     SplunkHecSourceResourceSchema,
	}
}

func NewLogStashSourceResource() resource.Resource {
	return &SourceResource[LogStashSourceModel]{
		typeName:          "logstash",
		fromModelFunc:     LogStashSourceFromModel,
		toModelFunc:       LogStashSourceToModel,
		getIdFunc:         func(m *LogStashSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LogStashSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     LogStashSourceResourceSchema,
	}
}

func NewFluentSourceResource() resource.Resource {
	return &SourceResource[FluentSourceModel]{
		typeName:          "fluent",
		fromModelFunc:     FluentSourceFromModel,
		toModelFunc:       FluentSourceToModel,
		getIdFunc:         func(m *FluentSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FluentSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     FluentSourceResourceSchema,
	}
}

func NewAzureEventHubSourceResource() resource.Resource {
	return &SourceResource[AzureEventHubSourceModel]{
		typeName:          "azure_event_hub",
		fromModelFunc:     AzureEventHubSourceFromModel,
		toModelFunc:       AzureEventHubSourceToModel,
		getIdFunc:         func(m *AzureEventHubSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AzureEventHubSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     AzureEventHubSourceResourceSchema,
	}
}

func NewKinesisFirehoseSourceResource() resource.Resource {
	return &SourceResource[KinesisFirehoseSourceModel]{
		typeName:          "kinesis_firehose",
		fromModelFunc:     KinesisFirehoseSourceFromModel,
		toModelFunc:       KinesisFirehoseSourceToModel,
		getIdFunc:         func(m *KinesisFirehoseSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KinesisFirehoseSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     KinesisFirehoseSourceResourceSchema,
	}
}

func NewLogAnalysisSourceResource() resource.Resource {
	return &SourceResource[LogAnalysisSourceModel]{
		typeName:          "log_analysis",
		fromModelFunc:     LogAnalysisSourceFromModel,
		toModelFunc:       LogAnalysisSourceToModel,
		getIdFunc:         func(m *LogAnalysisSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LogAnalysisSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     LogAnalysisSourceResourceSchema,
	}
}

func NewOpenTelemetryTracesSourceResource() resource.Resource {
	return &SourceResource[OpenTelemetryTracesSourceModel]{
		typeName:          "open_telemetry_traces",
		fromModelFunc:     OpenTelemetryTracesSourceFromModel,
		toModelFunc:       OpenTelemetryTracesSourceToModel,
		getIdFunc:         func(m *OpenTelemetryTracesSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *OpenTelemetryTracesSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     OpenTelemetryTracesSourceResourceSchema,
	}
}
