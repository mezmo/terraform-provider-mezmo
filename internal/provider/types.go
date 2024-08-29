package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
)

// Generic type representing a source / processor / destination model.
type ComponentModel interface {
	SourceModel | ProcessorModel | DestinationModel | AlertModel
}

type idGetterFunc[T ComponentModel] func(*T) basetypes.StringValue
type getSchemaFunc func() schema.Schema

type sourceToModelFunc[T ComponentModel] func(model *T, component *Source)
type sourceFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Source, diag.Diagnostics)

type processorToModelFunc[T ComponentModel] func(model *T, component *Processor)
type processorFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Processor, diag.Diagnostics)

type destinationToModelFunc[T ComponentModel] func(model *T, component *Destination)
type destinationFromModelFunc[T ComponentModel] func(model *T, previousState *T) (*Destination, diag.Diagnostics)

type alertToModelFunc[T AlertModel] func(model *T, component *Alert)
type alertFromModelFunc[T AlertModel] func(model *T, previousState *T) (*Alert, diag.Diagnostics)
