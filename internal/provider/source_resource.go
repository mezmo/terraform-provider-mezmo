package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/sources"
)

var (
	_ resource.Resource              = &PipelineResource{}
	_ resource.ResourceWithConfigure = &PipelineResource{}
)

func NewSourceResource() resource.Resource {
	return &SourceResource{}
}

type SourceResource struct {
	client client.Client
}

// Configure implements resource.ResourceWithConfigure.
func (r *SourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// TODO: Extract
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(client.Client)

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
func (r *SourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DemoSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	component := DemoSourceFromModel(&plan)
	stored, err := r.client.CreateSource(component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	DemoSourceToModel(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *SourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DemoSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteSource(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not source, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (*SourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_demo" + "_source"
}

// Read implements resource.Resource.
func (r *SourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DemoSourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Source(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			"Could not read source with id  "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	DemoSourceToModel(&state, component)
}

// Update implements resource.Resource.
func (r *SourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DemoSourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	var state DemoSourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component := DemoSourceFromModel(&plan)
	// Set id from the current state (not in plan)
	component.Id = state.Id.ValueString()
	stored, err := r.client.UpdateSource(component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			"Could not updated source, unexpected error: "+err.Error(),
		)
		return
	}

	DemoSourceToModel(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (*SourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DemoSourceResourceSchema()
}
