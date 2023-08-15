package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/transforms"
)

type TransformModel interface {
	StringifyTransformModel
}

type TransformResource[T TransformModel] struct {
	client            Client
	typeName          string
	fromModelFunc     componentFromModelFunc[T]
	toModelFunc       componentToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	getSchemaFunc     getSchemaFunc
}

// Configure implements resource.ResourceWithConfigure.
func (r *TransformResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Transform Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to Mezmo.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *TransformResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan T
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, dd := r.fromModelFunc(&plan, nil)
	if setDiagnosticsHasError(dd, &resp.Diagnostics) {
		return
	}
	stored, err := r.client.CreateTransform(r.getPipelineIdFunc(&plan).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating transform",
			"Could not create transform, unexpected error: "+err.Error(),
		)
		return
	}

	r.toModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *TransformResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteTransform(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting transform",
			"Could not delete transform, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *TransformResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName + "_transform"
}

// Read implements resource.Resource.
func (r *TransformResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Transform(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading transform",
			fmt.Sprintf("Could not read transform with id %s and pipeline_id %s: %s",
				r.getIdFunc(&state), r.getPipelineIdFunc(&state), err.Error()),
		)
		return
	}

	r.toModelFunc(&state, component)
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *TransformResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan T
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, dd := r.fromModelFunc(&plan, &state)
	if setDiagnosticsHasError(dd, &resp.Diagnostics) {
		return
	}
	stored, err := r.client.UpdateTransform(r.getPipelineIdFunc(&state).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating transform",
			"Could not update transform, unexpected error: "+err.Error(),
		)
		return
	}

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *TransformResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.getSchemaFunc()
}
