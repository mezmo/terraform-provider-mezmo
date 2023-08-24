package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/sinks"
)

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
