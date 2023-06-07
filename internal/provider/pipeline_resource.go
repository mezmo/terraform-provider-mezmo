package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/types"
)

var (
	_ resource.Resource              = &PipelineResource{}
	_ resource.ResourceWithConfigure = &PipelineResource{}
)

func NewPipelineResource() resource.Resource {
	return &PipelineResource{}
}

type PipelineResource struct {
	client client.Client
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
	var plan pipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline := Pipeline{
		Title: plan.Title.String(),
	}
	stored, err := r.client.CreatePipeline(&pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}
	plan.Id = StringValue(stored.Id)
	plan.UpdatedAt = StringValue(stored.UpdatedAt.Format(time.RFC3339))

	// Set the plan
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *PipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pipelineResourceModel
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
func (*PipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

// Read implements resource.Resource.
func (r *PipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	state.Title = StringValue(pipeline.Title)
	state.UpdatedAt = StringValue(pipeline.UpdatedAt.String())
}

// Schema implements resource.Resource.
func (*PipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = pipelineResourceSchema()
}

// Update implements resource.Resource.
func (r *PipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedAt, _ := time.Parse(time.RFC3339, plan.UpdatedAt.ValueString())

	pipeline := Pipeline{
		Id:        plan.Id.ValueString(),
		Title:     plan.Title.ValueString(),
		UpdatedAt: updatedAt,
	}
	stored, err := r.client.UpdatePipeline(&pipeline)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Pipeline",
			"Could not pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Title = StringValue(stored.Title)
	plan.UpdatedAt = StringValue(stored.UpdatedAt.String())
}
