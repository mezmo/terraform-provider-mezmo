package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models"
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

	// The "trick" here is that we don't actually want to change any of the
	// values based on the response. The response here will have the most
	// recent `updated_at` timestamp, but we need it to stay consistent
	// with the state of the pipeline's `updated_at` timestamp. This way,
	// we never have the plan values deviate from what's stored in the DB,
	// and the parent pipeline can trigger another publish when it refreshes
	// and sees that the `updated_at` timestamp has changed on the server.
	_, err := r.client.PublishPipeline(publish.PipelineId, ctx)

	if err != nil && err.(client.ApiResponseError).Code != "ENOCHANGES" {
		resp.Diagnostics.AddError(
			"Error publishing pipeline",
			"Could not publish pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the state to match the values from the parent pipeline
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
	// NOT having a `Read` for this resource is imperative to how it functions.
	// We always want the parent pipeline to control the changing of this resource
	// via `RequiresReplace`, so we never want to refresh the schema here.
}

func (*PublishPipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = PublishPipelineResourceSchema()
}

func (r *PublishPipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// There is no update for this resource since changed values will trigger a destroy/create.
}
