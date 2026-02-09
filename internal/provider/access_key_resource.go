package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models"
)

var (
	_ resource.Resource              = &AccessKeyResource{}
	_ resource.ResourceWithConfigure = &AccessKeyResource{}
)

func NewAccessKeyResource() resource.Resource {
	return &AccessKeyResource{}
}

type AccessKeyResource struct {
	client client.Client
}

func (r *AccessKeyResource) TypeName() string {
	return PROVIDER_TYPE_NAME + "_access_key"
}

func (r *AccessKeyResource) NodeType() string {
	return "accessKey"
}

func (r *AccessKeyResource) TerraformSchema() schema.Schema {
	return AccessKeyResourceSchema()
}

func (r *AccessKeyResource) NotConvertible() bool {
	// There is no import for access keys, so implement the "not convertible" interface
	return true
}

// Configure implements resource.ResourceWithConfigure.
func (r *AccessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// POST access key
func (r *AccessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AccessKeyResourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	accessKey := AccessKeyFromModel(&plan)

	stored, err := r.client.CreateAccessKey(accessKey, ctx)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access key",
			"Could not create access key, unexpected error: "+err.Error(),
		)
		return
	}

	AccessKeyToModel(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	setDiagnosticsHasError(diags, &resp.Diagnostics)
}

// Delete implements resource.Resource.
func (r *AccessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AccessKeyResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	accessKey := AccessKeyFromModel(&state)

	err := r.client.DeleteAccessKey(accessKey, ctx)
	if client.IsNotFoundError(err) {
		// If the key wasn't found, just ignore and let TF clean up state on its own
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Access Key",
			"Could not delete access key, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *AccessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

func (r *AccessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// There is no getter for access keys. Implement a noop to satisfy the interface.
}

func (*AccessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AccessKeyResourceSchema()
}

func (r *AccessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Updating keys is not supported. Throw an error until such time (if ever) we allow patching.
	resp.Diagnostics.AddError(
		fmt.Sprintf("Cannot update %s after creation", r.TypeName()),
		"Access keys are currently imutable and cannot be updated.",
	)
}
