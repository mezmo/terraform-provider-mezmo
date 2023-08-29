package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/sinks"
)

func NewAzureBlobStorageSinkResource() resource.Resource {
	return &SinkResource[AzureBlobStorageSinkModel]{
		typeName:          "azure_blob_storage",
		fromModelFunc:     AzureBlobStorageFromModel,
		toModelFunc:       AzureBlobStorageToModel,
		getIdFunc:         func(m *AzureBlobStorageSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *AzureBlobStorageSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     AzureBlobStorageResourceSchema,
	}
}

func NewBlackholeSinkResource() resource.Resource {
	return &SinkResource[BlackholeSinkModel]{
		typeName:          "blackhole",
		fromModelFunc:     BlackholeSinkFromModel,
		toModelFunc:       BlackholeSinkToModel,
		getIdFunc:         func(m *BlackholeSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *BlackholeSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     BlackholeSinkResourceSchema,
	}
}

func NewDatadogLogsSinkResource() resource.Resource {
	return &SinkResource[DatadogLogsSinkModel]{
		typeName:          "datadog-logs",
		fromModelFunc:     DatadogLogsFromModel,
		toModelFunc:       DatadogLogsSinkToModel,
		getIdFunc:         func(m *DatadogLogsSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogLogsSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     DatadogLogsSinkResourceSchema,
	}
}

func NewDatadogMetricsSinkResource() resource.Resource {
	return &SinkResource[DatadogMetricsSinkModel]{
		typeName:          "datadog-metrics",
		fromModelFunc:     DatadogMetricsFromModel,
		toModelFunc:       DatadogMetricsSinkToModel,
		getIdFunc:         func(m *DatadogMetricsSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *DatadogMetricsSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     DatadogMetricsSinkResourceSchema,
	}
}

func NewElasticSearchSinkResource() resource.Resource {
	return &SinkResource[ElasticSearchSinkModel]{
		typeName:          "elasticsearch",
		fromModelFunc:     ElasticSearchSinkFromModel,
		toModelFunc:       ElasticSearchSinkToModel,
		getIdFunc:         func(m *ElasticSearchSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *ElasticSearchSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     ElasticSearchSinkResourceSchema,
	}
}

func NewHoneycombLogsSinkResource() resource.Resource {
	return &SinkResource[HoneycombLogsSinkModel]{
		typeName:          "honeycomb_logs",
		fromModelFunc:     HoneycombLogsFromModel,
		toModelFunc:       HoneycombLogsToModel,
		getIdFunc:         func(m *HoneycombLogsSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HoneycombLogsSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     HoneycombLogsResourceSchema,
	}
}

func NewHttpSinkResource() resource.Resource {
	return &SinkResource[HttpSinkModel]{
		typeName:          "http",
		fromModelFunc:     HttpSinkFromModel,
		toModelFunc:       HttpSinkToModel,
		getIdFunc:         func(m *HttpSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *HttpSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     HttpSinkResourceSchema,
	}
}

func NewMezmoSinkResource() resource.Resource {
	return &SinkResource[MezmoSinkModel]{
		typeName:          "logs",
		fromModelFunc:     MezmoSinkFromModel,
		toModelFunc:       MezmoSinkToModel,
		getIdFunc:         func(m *MezmoSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *MezmoSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     MezmoSinkResourceSchema,
	}
}

func NewNewRelicSinkResource() resource.Resource {
	return &SinkResource[NewRelicSinkModel]{
		typeName:          "new_relic",
		fromModelFunc:     NewRelicSinkFromModel,
		toModelFunc:       NewRelicSinkToModel,
		getIdFunc:         func(m *NewRelicSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *NewRelicSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     NewRelicSinkResourceSchema,
	}
}

func NewPrometheusRemoteWriteSinkResource() resource.Resource {
	return &SinkResource[PrometheusRemoteWriteSinkModel]{
		typeName:          "prometheus_remote_write",
		fromModelFunc:     PrometheusRemoteWriteSinkFromModel,
		toModelFunc:       PrometheusRemoteWriteSinkToModel,
		getIdFunc:         func(m *PrometheusRemoteWriteSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *PrometheusRemoteWriteSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     PrometheusRemoteWriteSinkResourceSchema,
	}
}

func NewSplunkHecLogsSinkResource() resource.Resource {
	return &SinkResource[SplunkHecLogsSinkModel]{
		typeName:          "splunk_hec_logs",
		fromModelFunc:     SplunkHecLogsSinkFromModel,
		toModelFunc:       SplunkHecLogsSinkToModel,
		getIdFunc:         func(m *SplunkHecLogsSinkModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc: func(m *SplunkHecLogsSinkModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:     SplunkHecLogsSinkResourceSchema,
	}
}
