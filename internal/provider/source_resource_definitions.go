package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/sources"
)

func NewDemoSourceResource() resource.Resource {
	return &SourceResource[DemoSourceModel]{
		typeName:            "demo",
		sourceFromModelFunc: DemoSourceFromModel,
		sourceToModelFunc:   DemoSourceToModel,
		getIdFunc:           func(m *DemoSourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc:   func(m *DemoSourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:       DemoSourceResourceSchema,
	}
}

func NewS3SourceResource() resource.Resource {
	return &SourceResource[S3SourceModel]{
		typeName:            "s3",
		sourceFromModelFunc: S3SourceFromModel,
		sourceToModelFunc:   S3SourceToModel,
		getIdFunc:           func(m *S3SourceModel) basetypes.StringValue { return m.Id },
		getPipelineIdFunc:   func(m *S3SourceModel) basetypes.StringValue { return m.PipelineId },
		getSchemaFunc:       S3SourceResourceSchema,
	}
}
