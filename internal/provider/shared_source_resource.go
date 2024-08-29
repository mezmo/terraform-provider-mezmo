package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models"
	. "github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

var (
	_ resource.Resource                = &SharedSourceResource{}
	_ resource.ResourceWithConfigure   = &SharedSourceResource{}
	_ resource.ResourceWithImportState = &SharedSourceResource{}
)

func NewSharedSourceResource() resource.Resource {
	return &SharedSourceResource{}
}

type SharedSourceResource struct {
	client client.Client
}

func (r *SharedSourceResource) TypeName() string {
	return PROVIDER_TYPE_NAME + "_shared_source"
}

func (r *SharedSourceResource) NodeType() string {
	return "shared_source"
}

func (r *SharedSourceResource) TerraformSchema() schema.Schema {
	return SharedSourceResourceSchema()
}

func (*SharedSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SharedSourceResourceSchema()
}

func (r *SharedSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

func (r *SharedSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SharedSourceResource) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("component Value does not contain an interfaceable type")
	}

	source, ok := component.Interface().(client.SharedSource)
	if !ok {
		return nil, errors.New("component Value cannot be cast to a Shared Source")
	}

	var model SharedSourceResourceModel
	SharedSourceToModel(&model, &source)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *SharedSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// POST shared source
func (r *SharedSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SharedSourceResourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	source := SharedSourceFromModel(&plan)

	if os.Getenv("DEBUG_SHARED_SOURCE") == "1" {
		fmt.Println(Json("----- Shared Source TO Create api ---", source))
	}
	stored, err := r.client.CreateSharedSource(source, ctx)

	if os.Getenv("DEBUG_SHARED_SOURCE") == "1" {
		fmt.Println(Json("----- Shared Source FROM Create api ---", stored))
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating shared source",
			"Could not create shared source, unexpected error: "+err.Error(),
		)
		return
	}

	SharedSourceToModel(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	setDiagnosticsHasError(diags, &resp.Diagnostics)
}

func (r *SharedSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state SharedSourceResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	source, err := r.client.SharedSource(state.Id.ValueString(), ctx)
	// force re-creation of manually deleted resources
	if client.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SharedSource",
			"Could not read shared source with id  "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	SharedSourceToModel(&state, source)
}

func (r *SharedSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan SharedSourceResourceModel
	if diags := req.Plan.Get(ctx, &plan); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}
	var state SharedSourceResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	source := SharedSourceFromModel(&plan)
	// Set id from the current state (not in plan)
	source.Id = state.Id.ValueString()

	if os.Getenv("DEBUG_SHARED_SOURCE") == "1" {
		fmt.Println(Json("----- Shared Source TO Update api ---", source))
	}

	stored, err := r.client.UpdateSharedSource(source, ctx)

	if os.Getenv("DEBUG_SHARED_SOURCE") == "1" {
		fmt.Println(Json("----- Shared Source FROM Update api ---", stored))
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Shared Source",
			"Could not update shared source, unexpected error: "+err.Error(),
		)
		return
	}

	SharedSourceToModel(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *SharedSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SharedSourceResourceModel
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	source := SharedSourceFromModel(&state)

	err := r.client.DeleteSharedSource(source, ctx)
	if client.IsNotFoundError(err) {
		// If not found, just ignore and let TF clean up state on its own
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Shared Source",
			"Could not delete shared source, unexpected error: "+err.Error(),
		)
		return
	}
}
