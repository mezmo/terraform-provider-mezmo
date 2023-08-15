package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
)

// Generic type representing a source / transform / sink model.
type ComponentModel interface {
	SourceModel | TransformModel | SinkModel
}

type idGetterFunc[T ComponentModel] func(*T) basetypes.StringValue
type componentToModelFunc[T ComponentModel] func(model *T, component *Component)
type componentFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Component, diag.Diagnostics)
type getSchemaFunc func() schema.Schema
