package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models"
)

var (
	_ resource.Resource              = &PublishPipelineResource{}
	_ resource.ResourceWithConfigure = &PublishPipelineResource{}
)

func NewPublishPipelineResource() resource.Resource {
	return &PublishPipelineResource{}
}

type PublishPipelineResource struct {
	client client.Client
}

func (r *PublishPipelineResource) TypeName() string {
	return PROVIDER_TYPE_NAME + "_publish_pipeline"
}

func (r *PublishPipelineResource) NodeType() string {
	return "publish_pipeline"
}

func (r *PublishPipelineResource) TerraformSchema() schema.Schema {
	return PublishPipelineResourceSchema()
}

func (r *PublishPipelineResource) NotConvertible() bool {
	// FIXME: Import is not implemented yet. We need LOG-20061 so that import/export can use child modules.
	return true
}

func (r *PublishPipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// POST pipeline publish
func (r *PublishPipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PublishPipelineResourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	publish := PublishPipelineFromModel(&plan)
	// Result doesn't matter, in fact it would cause data inconsistencies if we used it.
	_, err := r.client.PublishPipeline(publish.PipelineId, ctx)

	if err != nil && err.(client.ApiResponseError).Code != "ENOCHANGES" {
		resp.Diagnostics.AddError(
			"Error publishing pipeline",
			"Could not publish pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the state to satisfy the requirement that the response state be set in this method.
	// The value of state is not used because it's removed on every `Read()`, but we need to prevent
	// data inconsistecies between the plan and result. There's no need to instantiate
	// a model here because it's only `pipeline_id`, so we can just use the `plan` as-is.
	diags := resp.State.Set(ctx, plan)
	setDiagnosticsHasError(diags, &resp.Diagnostics)
}

func (r *PublishPipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// There is no action when this resource is deleted.
}

func (r *PublishPipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

func (r *PublishPipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is the magic. We *always* want the publish to fire, so force it to be re-created on every plan
	// by removing it from state. Without doing this, the resource might not be selected for a "change" and would
	// skip publishing most times.
	resp.State.RemoveResource(ctx)
}

func (*PublishPipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PublishPipelineResourceSchema()
}

func (r *PublishPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// There is no update for this resource since changed values will trigger a destroy/create.
}
