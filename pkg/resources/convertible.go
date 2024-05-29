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
type AlertApiModel = client.Alert
type BaseApiModel = client.BaseNode
type SharedSourceApiModel = client.SharedSource

type ConvertibleResourceDef interface {
	/// the terraform resource type. example: mezmo_http_source
	TypeName() string
	/// the pipeline service node type. example: demo_logs_source, threshold_alert
	NodeType() string
	TerraformSchema() schema.Schema
	ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error)
}

type NotConvertibleResourceDef interface {
	NotConvertible() bool
}

func ConvertibleResources() ([]ConvertibleResourceDef, error) {
	p := provider.MezmoProvider{}
	terraformRes := p.Resources(context.Background())
	convertibleRes := []ConvertibleResourceDef{}
	for _, tfResFn := range terraformRes {
		tfRes := tfResFn()
		cRes, ok := tfRes.(ConvertibleResourceDef)
		if !ok {
			not_convertible, ok := tfRes.(NotConvertibleResourceDef)
			if !ok {
				err := fmt.Errorf("terraform resource %T does not implement ConvertibleResourceDef nor NotConvertibleResourcedef", tfRes)
				return nil, err
			}
			if !not_convertible.NotConvertible() {
				err := fmt.Errorf("terraform resource %T must return `true` from `NotConvertible()`", tfRes)
				return nil, err
			}
			continue
		}
		convertibleRes = append(convertibleRes, cRes)
	}
	return convertibleRes, nil
}
