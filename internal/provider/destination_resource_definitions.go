package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/destinations"
)

func NewAzureBlobStorageDestinationResource() resource.Resource {
	return &DestinationResource[AzureBlobStorageDestinationModel]{
		typeName:          AZURE_BLOB_STORAGE_DESTINATION_TYPE_NAME,
		nodeName:          AZURE_BLOB_STORAGE_DESTINATION_NODE_NAME,
		fromModelFunc:     AzureBlobStorageFromModel,
		toModelFunc:       AzureBlobStorageToModel,
		getIdFunc:         func(m *AzureBlobStorageDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AzureBlobStorageDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            AzureBlobStorageResourceSchema,
	}
}

func NewBlackholeDestinationResource() resource.Resource {
	return &DestinationResource[BlackholeDestinationModel]{
		typeName:          BLACKHOLE_DESTINATION_TYPE_NAME,
		nodeName:          BLACKHOLE_DESTINATION_NODE_NAME,
		fromModelFunc:     BlackholeDestinationFromModel,
		toModelFunc:       BlackholeDestinationToModel,
		getIdFunc:         func(m *BlackholeDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *BlackholeDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            BlackholeDestinationResourceSchema,
	}
}

func NewDatadogLogsDestinationResource() resource.Resource {
	return &DestinationResource[DatadogLogsDestinationModel]{
		typeName:          DATADOG_LOGS_DESTINATION_TYPE_NAME,
		nodeName:          DATADOG_LOGS_DESTINATION_NODE_NAME,
		fromModelFunc:     DatadogLogsFromModel,
		toModelFunc:       DatadogLogsDestinationToModel,
		getIdFunc:         func(m *DatadogLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            DatadogLogsDestinationResourceSchema,
	}
}

func NewDatadogMetricsDestinationResource() resource.Resource {
	return &DestinationResource[DatadogMetricsDestinationModel]{
		typeName:          DATADOG_METRICS_DESTINATION_TYPE_NAME,
		nodeName:          DATADOG_METRICS_DESTINATION_NODE_NAME,
		fromModelFunc:     DatadogMetricsFromModel,
		toModelFunc:       DatadogMetricsDestinationToModel,
		getIdFunc:         func(m *DatadogMetricsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogMetricsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            DatadogMetricsDestinationResourceSchema,
	}
}

func NewElasticSearchDestinationResource() resource.Resource {
	return &DestinationResource[ElasticSearchDestinationModel]{
		typeName:          ELASTICSEARCH_DESTINATION_TYPE_NAME,
		nodeName:          ELASTICSEARCH_DESTINATION_NODE_NAME,
		fromModelFunc:     ElasticSearchDestinationFromModel,
		toModelFunc:       ElasticSearchDestinationToModel,
		getIdFunc:         func(m *ElasticSearchDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ElasticSearchDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            ElasticSearchDestinationResourceSchema,
	}
}

func NewHoneycombLogsDestinationResource() resource.Resource {
	return &DestinationResource[HoneycombLogsDestinationModel]{
		typeName:          HONEYCOMB_LOGS_DESTINATION_TYPE_NAME,
		nodeName:          HONEYCOMB_LOGS_DESTINATION_NODE_NAME,
		fromModelFunc:     HoneycombLogsFromModel,
		toModelFunc:       HoneycombLogsToModel,
		getIdFunc:         func(m *HoneycombLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HoneycombLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            HoneycombLogsResourceSchema,
	}
}

func NewHttpDestinationResource() resource.Resource {
	return &DestinationResource[HttpDestinationModel]{
		typeName:          HTTP_DESTINATION_TYPE_NAME,
		nodeName:          HTTP_DESTINATION_NODE_NAME,
		fromModelFunc:     HttpDestinationFromModel,
		toModelFunc:       HttpDestinationToModel,
		getIdFunc:         func(m *HttpDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HttpDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            HttpDestinationResourceSchema,
	}
}

func NewKafkaDestinationResource() resource.Resource {
	return &DestinationResource[KafkaDestinationModel]{
		typeName:          KAFKA_DESTINATION_TYPE_NAME,
		nodeName:          KAFKA_DESTINATION_NODE_NAME,
		fromModelFunc:     KafkaDestinationFromModel,
		toModelFunc:       KafkaDestinationToModel,
		getIdFunc:         func(m *KafkaDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KafkaDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            KafkaDestinationResourceSchema,
	}
}

func NewLokiDestinationResource() resource.Resource {
	return &DestinationResource[LokiDestinationModel]{
		typeName:          LOKI_DESTINATION_TYPE_NAME,
		nodeName:          LOKI_DESTINATION_NODE_NAME,
		fromModelFunc:     LokiFromModel,
		toModelFunc:       LokiDestinationToModel,
		getIdFunc:         func(m *LokiDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LokiDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            LokiDestinationResourceSchema,
	}
}

func NewMezmoDestinationResource() resource.Resource {
	return &DestinationResource[MezmoDestinationModel]{
		typeName:          MEZMO_DESTINATION_TYPE_NAME,
		nodeName:          MEZMO_DESTINATION_NODE_NAME,
		fromModelFunc:     MezmoDestinationFromModel,
		toModelFunc:       MezmoDestinationToModel,
		getIdFunc:         func(m *MezmoDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MezmoDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            MezmoDestinationResourceSchema,
	}
}

func NewNewRelicDestinationResource() resource.Resource {
	return &DestinationResource[NewRelicDestinationModel]{
		typeName:          NEWRELIC_DESTINATION_TYPE_NAME,
		nodeName:          NEWRELIC_DESTINATION_NODE_NAME,
		fromModelFunc:     NewRelicDestinationFromModel,
		toModelFunc:       NewRelicDestinationToModel,
		getIdFunc:         func(m *NewRelicDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *NewRelicDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            NewRelicDestinationResourceSchema,
	}
}

func NewPrometheusRemoteWriteDestinationResource() resource.Resource {
	return &DestinationResource[PrometheusRemoteWriteDestinationModel]{
		typeName:          PROMETHEUS_REMOTE_WRITE_DESTINATION_TYPE_NAME,
		nodeName:          PROMETHEUS_REMOTE_WRITE_DESTINATION_NODE_NAME,
		fromModelFunc:     PrometheusRemoteWriteDestinationFromModel,
		toModelFunc:       PrometheusRemoteWriteDestinationToModel,
		getIdFunc:         func(m *PrometheusRemoteWriteDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *PrometheusRemoteWriteDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            PrometheusRemoteWriteDestinationResourceSchema,
	}
}

func NewS3DestinationResource() resource.Resource {
	return &DestinationResource[S3DestinationModel]{
		typeName:          S3_DESTINATION_TYPE_NAME,
		nodeName:          S3_DESTINATION_NODE_NAME,
		fromModelFunc:     S3DestinationFromModel,
		toModelFunc:       S3DestinationToModel,
		getIdFunc:         func(m *S3DestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *S3DestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            S3DestinationResourceSchema,
	}
}

func NewSplunkHecLogsDestinationResource() resource.Resource {
	return &DestinationResource[SplunkHecLogsDestinationModel]{
		typeName:          SPLUNK_HEC_LOGS_DESTINATION_TYPE_NAME,
		nodeName:          SPLUNK_HEC_LOGS_DESTINATION_NODE_NAME,
		fromModelFunc:     SplunkHecLogsDestinationFromModel,
		toModelFunc:       SplunkHecLogsDestinationToModel,
		getIdFunc:         func(m *SplunkHecLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SplunkHecLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            SplunkHecLogsDestinationResourceSchema,
	}
}

func NewGcpCloudStorageDestinationResource() resource.Resource {
	return &DestinationResource[GcpCloudStorageDestinationModel]{
		typeName:          GCP_CLOUD_STORAGE_DESTINATION_TYPE_NAME,
		nodeName:          GCP_CLOUD_STORAGE_DESTINATION_NODE_NAME,
		fromModelFunc:     GcpCloudStorageDestinationFromModel,
		toModelFunc:       GcpCloudStorageDestinationToModel,
		getIdFunc:         func(m *GcpCloudStorageDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *GcpCloudStorageDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            GcpCloudStorageResourceSchema,
	}
}

func NewGcpCloudMonitoringDestinationResource() resource.Resource {
	return &DestinationResource[GcpCloudMonitoringDestinationModel]{
		typeName:          GCP_CLOUD_MONITORING_DESTINATION_TYPE_NAME,
		nodeName:          GCP_CLOUD_MONITORING_DESTINATION_NODE_NAME,
		fromModelFunc:     GcpCloudMonitoringDestinationFromModel,
		toModelFunc:       GcpCloudMonitoringDestinationToModel,
		getIdFunc:         func(m *GcpCloudMonitoringDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *GcpCloudMonitoringDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            GcpCloudMonitoringResourceSchema,
	}
}

func NewGcpCloudOperationsDestinationResource() resource.Resource {
	return &DestinationResource[GcpCloudOperationsDestinationModel]{
		typeName:          GCP_CLOUD_OPERATIONS_DESTINATION_TYPE_NAME,
		nodeName:          GCP_CLOUD_OPERATIONS_DESTINATION_NODE_NAME,
		fromModelFunc:     GcpCloudOperationsDestinationFromModel,
		toModelFunc:       GcpCloudOperationsDestinationToModel,
		getIdFunc:         func(m *GcpCloudOperationsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *GcpCloudOperationsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            GcpCloudOperationsResourceSchema,
	}
}

func NewGcpCloudPubSubDestinationResource() resource.Resource {
	return &DestinationResource[GcpCloudPubSubDestinationModel]{
		typeName:          GCP_CLOUD_PUBSUB_DESTINATION_TYPE_NAME,
		nodeName:          GCP_CLOUD_PUBSUB_DESTINATION_NODE_NAME,
		fromModelFunc:     GcpCloudPubSubDestinationFromModel,
		toModelFunc:       GcpCloudPubSubDestinationToModel,
		getIdFunc:         func(m *GcpCloudPubSubDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *GcpCloudPubSubDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            GcpCloudPubSubResourceSchema,
	}
}
