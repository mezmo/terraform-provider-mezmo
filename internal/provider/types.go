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
type getSchemaFunc func() schema.Schema

type sourceToModelFunc[T ComponentModel] func(model *T, component *Source)
type sourceFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Source, diag.Diagnostics)

type transformToModelFunc[T ComponentModel] func(model *T, component *Transform)
type transformFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Transform, diag.Diagnostics)

type sinkToModelFunc[T ComponentModel] func(model *T, component *Sink)
type sinkFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Sink, diag.Diagnostics)
