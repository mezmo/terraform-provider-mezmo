package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/sources"
)

func NewDemoSourceResource() resource.Resource {
	return &SourceResource[DemoSourceModel]{
		typeName:          DEMO_SOURCE_TYPE_NAME,
		nodeName:          DEMO_SOURCE_NODE_NAME,
		fromModelFunc:     DemoSourceFromModel,
		toModelFunc:       DemoSourceToModel,
		getIdFunc:         func(m *DemoSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DemoSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            DemoSourceResourceSchema,
	}
}

func NewAgentSourceResource() resource.Resource {
	return &SourceResource[AgentSourceModel]{
		typeName:          AGENT_SOURCE_TYPE_NAME,
		nodeName:          AGENT_SOURCE_NODE_NAME,
		fromModelFunc:     AgentSourceFromModel,
		toModelFunc:       AgentSourceToModel,
		getIdFunc:         func(m *AgentSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AgentSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            AgentSourceResourceSchema,
	}
}

func NewKafkaSourceResource() resource.Resource {
	return &SourceResource[KafkaSourceModel]{
		typeName:          KAFKA_SOURCE_TYPE_NAME,
		nodeName:          KAFKA_SOURCE_NODE_NAME,
		fromModelFunc:     KafkaSourceFromModel,
		toModelFunc:       KafkaSourceToModel,
		getIdFunc:         func(m *KafkaSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KafkaSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            KafkaSourceResourceSchema,
	}
}

func NewPrometheusRemoteWriteSourceResource() resource.Resource {
	return &SourceResource[PrometheusRemoteWriteSourceModel]{
		typeName:          PROMETHEUS_REMOTE_WRITE_SOURCE_TYPE_NAME,
		nodeName:          PROMETHEUS_REMOTE_WRITE_SOURCE_NODE_NAME,
		fromModelFunc:     PrometheusRemoteWriteSourceFromModel,
		toModelFunc:       PrometheusRemoteWriteSourceToModel,
		getIdFunc:         func(m *PrometheusRemoteWriteSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *PrometheusRemoteWriteSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            PrometheusRemoteWriteSourceResourceSchema,
	}
}

func NewS3SourceResource() resource.Resource {
	return &SourceResource[S3SourceModel]{
		typeName:          S3_SOURCE_TYPE_NAME,
		nodeName:          S3_SOURCE_NODE_NAME,
		fromModelFunc:     S3SourceFromModel,
		toModelFunc:       S3SourceToModel,
		getIdFunc:         func(m *S3SourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *S3SourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            S3SourceResourceSchema,
	}
}

func NewHttpSourceResource() resource.Resource {
	return &SourceResource[HttpSourceModel]{
		typeName:          HTTP_SOURCE_TYPE_NAME,
		nodeName:          HTTP_SOURCE_NODE_NAME,
		fromModelFunc:     HttpSourceFromModel,
		toModelFunc:       HttpSourceToModel,
		getIdFunc:         func(m *HttpSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HttpSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            HttpSourceResourceSchema,
	}
}

func NewSQSSourceResource() resource.Resource {
	return &SourceResource[SQSSourceModel]{
		typeName:          SQS_SOURCE_TYPE_NAME,
		nodeName:          SQS_SOURCE_NODE_NAME,
		fromModelFunc:     SQSSourceFromModel,
		toModelFunc:       SQSSourceToModel,
		getIdFunc:         func(m *SQSSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SQSSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            SQSSourceResourceSchema,
	}
}

func NewSplunkHecSourceResource() resource.Resource {
	return &SourceResource[SplunkHecSourceModel]{
		typeName:          SPLUNK_HEC_SOURCE_TYPE_NAME,
		nodeName:          SPLUNK_HEC_SOURCE_NODE_NAME,
		fromModelFunc:     SplunkHecSourceFromModel,
		toModelFunc:       SplunkHecSourceToModel,
		getIdFunc:         func(m *SplunkHecSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SplunkHecSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            SplunkHecSourceResourceSchema,
	}
}

func NewLogStashSourceResource() resource.Resource {
	return &SourceResource[LogStashSourceModel]{
		typeName:          LOGSTASH_SOURCE_TYPE_NAME,
		nodeName:          LOGSTASH_SOURCE_NODE_NAME,
		fromModelFunc:     LogStashSourceFromModel,
		toModelFunc:       LogStashSourceToModel,
		getIdFunc:         func(m *LogStashSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LogStashSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            LogStashSourceResourceSchema,
	}
}

func NewFluentSourceResource() resource.Resource {
	return &SourceResource[FluentSourceModel]{
		typeName:          FLUENT_SOURCE_TYPE_NAME,
		nodeName:          FLUENT_SOURCE_NODE_NAME,
		fromModelFunc:     FluentSourceFromModel,
		toModelFunc:       FluentSourceToModel,
		getIdFunc:         func(m *FluentSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *FluentSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            FluentSourceResourceSchema,
	}
}

func NewAzureEventHubSourceResource() resource.Resource {
	return &SourceResource[AzureEventHubSourceModel]{
		typeName:          AZURE_EVENT_HUB_SOURCE_TYPE_NAME,
		nodeName:          AZURE_EVENT_HUB_SOURCE_NODE_NAME,
		fromModelFunc:     AzureEventHubSourceFromModel,
		toModelFunc:       AzureEventHubSourceToModel,
		getIdFunc:         func(m *AzureEventHubSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AzureEventHubSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            AzureEventHubSourceResourceSchema,
	}
}

func NewKinesisFirehoseSourceResource() resource.Resource {
	return &SourceResource[KinesisFirehoseSourceModel]{
		typeName:          KINESIS_FIREHOSE_SOURCE_TYPE_NAME,
		nodeName:          KINESIS_FIREHOSE_SOURCE_NODE_NAME,
		fromModelFunc:     KinesisFirehoseSourceFromModel,
		toModelFunc:       KinesisFirehoseSourceToModel,
		getIdFunc:         func(m *KinesisFirehoseSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KinesisFirehoseSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            KinesisFirehoseSourceResourceSchema,
	}
}

func NewLogAnalysisSourceResource() resource.Resource {
	return &SourceResource[LogAnalysisSourceModel]{
		typeName:          LOG_ANALYSIS_SOURCE_TYPE_NAME,
		nodeName:          LOG_ANALYSIS_SOURCE_NODE_NAME,
		fromModelFunc:     LogAnalysisSourceFromModel,
		toModelFunc:       LogAnalysisSourceToModel,
		getIdFunc:         func(m *LogAnalysisSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LogAnalysisSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            LogAnalysisSourceResourceSchema,
	}
}

func NewWebhookSourceResource() resource.Resource {
	return &SourceResource[WebhookSourceModel]{
		typeName:          WEBHOOK_SOURCE_TYPE_NAME,
		nodeName:          WEBHOOK_SOURCE_NODE_NAME,
		fromModelFunc:     WebhookSourceFromModel,
		toModelFunc:       WebhookSourceToModel,
		getIdFunc:         func(m *WebhookSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *WebhookSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            WebhookSourceResourceSchema,
	}
}

func NewDatadogSourceResource() resource.Resource {
	return &SourceResource[DatadogSourceModel]{
		typeName:          DATADOG_SOURCE_TYPE_NAME,
		nodeName:          DATADOG_SOURCE_NODE_NAME,
		fromModelFunc:     DatadogSourceFromModel,
		toModelFunc:       DatadogSourceToModel,
		getIdFunc:         func(m *DatadogSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            DatadogSourceResourceSchema,
	}
}

func NewOpenTelemetryLogsSourceResource() resource.Resource {
	return &SourceResource[OpenTelemetryLogsSourceModel]{
		typeName:          OPEN_TELEMETRY_LOGS_SOURCE_TYPE_NAME,
		nodeName:          OPEN_TELEMETRY_LOGS_SOURCE_NODE_NAME,
		fromModelFunc:     OpenTelemetryLogsSourceFromModel,
		toModelFunc:       OpenTelemetryLogsSourceToModel,
		getIdFunc:         func(m *OpenTelemetryLogsSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *OpenTelemetryLogsSourceModel) basetypes.StringValue { return m.PipelineId },
		schema:            OpenTelemetryLogsSourceResourceSchema,
	}
}
