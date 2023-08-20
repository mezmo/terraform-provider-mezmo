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
