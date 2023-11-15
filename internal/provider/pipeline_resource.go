package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models"
)

var (
	_ resource.Resource                = &PipelineResource{}
	_ resource.ResourceWithConfigure   = &PipelineResource{}
	_ resource.ResourceWithImportState = &PipelineResource{}
)

func NewPipelineResource() resource.Resource {
	return &PipelineResource{}
}

type PipelineResource struct {
	client client.Client
}

func (r *PipelineResource) TypeName() string {
	return PROVIDER_TYPE_NAME + "_pipeline"
}

func (r *PipelineResource) NodeType() string {
	return "pipeline"
}

func (r *PipelineResource) TerraformSchema() schema.Schema {
	return PipelineResourceSchema()
}

func (r *PipelineResource) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("component Value does not contain an interfaceable type")
	}

	pipeline, ok := component.Interface().(client.Pipeline)
	if !ok {
		return nil, errors.New("component Value cannot be cast to Pipeline")
	}

	var model PipelineResourceModel
	PipelineToModel(&model, &pipeline)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *PipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan PipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline := PipelineFromModel(&plan)
	stored, err := r.client.CreatePipeline(pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	PipelineToModel(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *PipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeletePipeline(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Pipeline",
			"Could not pipeline, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *PipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

// Read implements resource.Resource.
func (r *PipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state PipelineResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	pipeline, err := r.client.Pipeline(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Pipeline",
			"Could not read pipeline with id  "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	PipelineToModel(&state, pipeline)
}

// Schema implements resource.Resource.
func (*PipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PipelineResourceSchema()
}

// Update implements resource.Resource.
func (r *PipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan PipelineResourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	var state PipelineResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	pipeline := PipelineFromModel(&plan)
	// Set id from the current state (not in plan)
	pipeline.Id = state.Id.ValueString()
	stored, err := r.client.UpdatePipeline(pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Pipeline",
			"Could not update pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	PipelineToModel(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func setDiagnosticsHasError(source diag.Diagnostics, target *diag.Diagnostics) bool {
	target.Append(source...)
	return target.HasError()
}

func (r *PipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
