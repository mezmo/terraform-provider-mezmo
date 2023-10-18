package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/destinations"
)

func NewAzureBlobStorageDestinationResource() resource.Resource {
	return &DestinationResource[AzureBlobStorageDestinationModel]{
		typeName:          "azure_blob_storage",
		fromModelFunc:     AzureBlobStorageFromModel,
		toModelFunc:       AzureBlobStorageToModel,
		getIdFunc:         func(m *AzureBlobStorageDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AzureBlobStorageDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            AzureBlobStorageResourceSchema,
	}
}

func NewBlackholeDestinationResource() resource.Resource {
	return &DestinationResource[BlackholeDestinationModel]{
		typeName:          "blackhole",
		fromModelFunc:     BlackholeDestinationFromModel,
		toModelFunc:       BlackholeDestinationToModel,
		getIdFunc:         func(m *BlackholeDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *BlackholeDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            BlackholeDestinationResourceSchema,
	}
}

func NewDatadogLogsDestinationResource() resource.Resource {
	return &DestinationResource[DatadogLogsDestinationModel]{
		typeName:          "datadog_logs",
		fromModelFunc:     DatadogLogsFromModel,
		toModelFunc:       DatadogLogsDestinationToModel,
		getIdFunc:         func(m *DatadogLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            DatadogLogsDestinationResourceSchema,
	}
}

func NewDatadogMetricsDestinationResource() resource.Resource {
	return &DestinationResource[DatadogMetricsDestinationModel]{
		typeName:          "datadog_metrics",
		fromModelFunc:     DatadogMetricsFromModel,
		toModelFunc:       DatadogMetricsDestinationToModel,
		getIdFunc:         func(m *DatadogMetricsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogMetricsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            DatadogMetricsDestinationResourceSchema,
	}
}

func NewElasticSearchDestinationResource() resource.Resource {
	return &DestinationResource[ElasticSearchDestinationModel]{
		typeName:          "elasticsearch",
		fromModelFunc:     ElasticSearchDestinationFromModel,
		toModelFunc:       ElasticSearchDestinationToModel,
		getIdFunc:         func(m *ElasticSearchDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ElasticSearchDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            ElasticSearchDestinationResourceSchema,
	}
}

func NewHoneycombLogsDestinationResource() resource.Resource {
	return &DestinationResource[HoneycombLogsDestinationModel]{
		typeName:          "honeycomb_logs",
		fromModelFunc:     HoneycombLogsFromModel,
		toModelFunc:       HoneycombLogsToModel,
		getIdFunc:         func(m *HoneycombLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HoneycombLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            HoneycombLogsResourceSchema,
	}
}

func NewHttpDestinationResource() resource.Resource {
	return &DestinationResource[HttpDestinationModel]{
		typeName:          "http",
		fromModelFunc:     HttpDestinationFromModel,
		toModelFunc:       HttpDestinationToModel,
		getIdFunc:         func(m *HttpDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HttpDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            HttpDestinationResourceSchema,
	}
}

func NewKafkaDestinationResource() resource.Resource {
	return &DestinationResource[KafkaDestinationModel]{
		typeName:          "kafka",
		fromModelFunc:     KafkaDestinationFromModel,
		toModelFunc:       KafkaDestinationToModel,
		getIdFunc:         func(m *KafkaDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *KafkaDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            KafkaDestinationResourceSchema,
	}
}

func NewLokiDestinationResource() resource.Resource {
	return &DestinationResource[LokiDestinationModel]{
		typeName:          "loki",
		fromModelFunc:     LokiFromModel,
		toModelFunc:       LokiDestinationToModel,
		getIdFunc:         func(m *LokiDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *LokiDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            LokiDestinationResourceSchema,
	}
}

func NewMezmoDestinationResource() resource.Resource {
	return &DestinationResource[MezmoDestinationModel]{
		typeName:          "logs",
		fromModelFunc:     MezmoDestinationFromModel,
		toModelFunc:       MezmoDestinationToModel,
		getIdFunc:         func(m *MezmoDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MezmoDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            MezmoDestinationResourceSchema,
	}
}

func NewNewRelicDestinationResource() resource.Resource {
	return &DestinationResource[NewRelicDestinationModel]{
		typeName:          "new_relic",
		fromModelFunc:     NewRelicDestinationFromModel,
		toModelFunc:       NewRelicDestinationToModel,
		getIdFunc:         func(m *NewRelicDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *NewRelicDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            NewRelicDestinationResourceSchema,
	}
}

func NewPrometheusRemoteWriteDestinationResource() resource.Resource {
	return &DestinationResource[PrometheusRemoteWriteDestinationModel]{
		typeName:          "prometheus_remote_write",
		fromModelFunc:     PrometheusRemoteWriteDestinationFromModel,
		toModelFunc:       PrometheusRemoteWriteDestinationToModel,
		getIdFunc:         func(m *PrometheusRemoteWriteDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *PrometheusRemoteWriteDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            PrometheusRemoteWriteDestinationResourceSchema,
	}
}

func NewS3DestinationResource() resource.Resource {
	return &DestinationResource[S3DestinationModel]{
		typeName:          "s3",
		fromModelFunc:     S3DestinationFromModel,
		toModelFunc:       S3DestinationToModel,
		getIdFunc:         func(m *S3DestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *S3DestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            S3DestinationResourceSchema,
	}
}

func NewSplunkHecLogsDestinationResource() resource.Resource {
	return &DestinationResource[SplunkHecLogsDestinationModel]{
		typeName:          "splunk_hec_logs",
		fromModelFunc:     SplunkHecLogsDestinationFromModel,
		toModelFunc:       SplunkHecLogsDestinationToModel,
		getIdFunc:         func(m *SplunkHecLogsDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SplunkHecLogsDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            SplunkHecLogsDestinationResourceSchema,
	}
}

func NewGcpCloudStorageDestinationResource() resource.Resource {
	return &DestinationResource[GcpCloudStorageDestinationModel]{
		typeName:          "gcp_cloud_storage",
		fromModelFunc:     GcpCloudStorageDestinationFromModel,
		toModelFunc:       GcpCloudStorageDestinationToModel,
		getIdFunc:         func(m *GcpCloudStorageDestinationModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *GcpCloudStorageDestinationModel) basetypes.StringValue { return m.PipelineId },
		schema:            GcpCloudStorageResourceSchema,
	}
}
