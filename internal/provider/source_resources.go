package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/sources"
)

type SourceModel interface {
	DemoSourceModel
}

type SourceResource[T SourceModel] struct {
	client              Client
	typeName            string
	sourceFromModelFunc func(model *T) *Component
	sourceToModelFunc   func(model *T, component *Component)
	getIdFunc           func(m *T) basetypes.StringValue
	getSchemaFunc       func() schema.Schema
}

func NewDemoSourceResource() resource.Resource {
	return &SourceResource[DemoSourceModel]{
		typeName:            "demo",
		sourceFromModelFunc: DemoSourceFromModel,
		sourceToModelFunc:   DemoSourceToModel,
		getIdFunc:           func(m *DemoSourceModel) basetypes.StringValue { return m.Id },
		getSchemaFunc:       DemoSourceResourceSchema,
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *SourceResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to Mezmo.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *SourceResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan T
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	component := r.sourceFromModelFunc(&plan)
	stored, err := r.client.CreateSource(component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	r.sourceToModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *SourceResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteSource(r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not source, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *SourceResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName + "_source"
}

// Read implements resource.Resource.
func (r *SourceResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Source(r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			"Could not read source with id  "+r.getIdFunc(&state).ValueString()+": "+err.Error(),
		)
		return
	}

	r.sourceToModelFunc(&state, component)
}

// Update implements resource.Resource.
func (r *SourceResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan T
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component := r.sourceFromModelFunc(&plan)
	// Set id from the current state (not in plan)
	component.Id = r.getIdFunc(&state).ValueString()
	stored, err := r.client.UpdateSource(component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			"Could not updated source, unexpected error: "+err.Error(),
		)
		return
	}

	r.sourceToModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *SourceResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.getSchemaFunc()
}
