package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/mezmo/terraform-provider-mezmo/internal/client"
	"github.com/mezmo/terraform-provider-mezmo/internal/provider"
	"reflect"
)

type PipelineApiModel = client.Pipeline
type SourceApiModel = client.Source
type ProcessorApiModel = client.Processor
type DestinationApiModel = client.Destination

type ConvertibleResourceDef interface {
	TypeName() string
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
			err := errors.New(fmt.Sprintf("terraform resource %T does not implement ConvertibleResourceDef", tfRes))
			return nil, err
		}
		convertibleRes[i] = cRes
	}
	return convertibleRes, nil
}
