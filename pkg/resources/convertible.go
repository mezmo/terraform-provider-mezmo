package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mezmo/terraform-provider-mezmo/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider"
)

type PipelineApiModel = client.Pipeline
type SourceApiModel = client.Source
type ProcessorApiModel = client.Processor
type DestinationApiModel = client.Destination
type BaseApiModel = client.BaseNode

type ConvertibleResourceDef interface {
	/// the terraform resource type. example: mezmo_http_source
	TypeName() string
	/// the pipeline service node type. example: demo_logs_source
	NodeType() string
	TerraformSchema() schema.Schema
	ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error)
}

func ConvertibleResources() ([]ConvertibleResourceDef, error) {
	p := provider.MezmoProvider{}
	terraformRes := p.Resources(context.Background())
	convertibleRes := make([]ConvertibleResourceDef, len(terraformRes))
	for i, tfResFn := range terraformRes {
		tfRes := tfResFn()
		cRes, ok := tfRes.(ConvertibleResourceDef)
		if !ok {
			err := fmt.Errorf("terraform resource %T does not implement ConvertibleResourceDef", tfRes)
			return nil, err
		}
		convertibleRes[i] = cRes
	}
	return convertibleRes, nil
}
