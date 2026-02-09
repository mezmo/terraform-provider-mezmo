package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/alerts"
)

type AlertModel interface {
	ThresholdAlertModel |
		ChangeAlertModel |
		AbsenceAlertModel
}

type AlertResource[T AlertModel] struct {
	client            Client
	typeName          string // The name to use as part of the terraform resource name: mezmo_{typeName}_alert
	fromModelFunc     alertFromModelFunc[T]
	toModelFunc       alertToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	schema            schema.Schema
}

func (r *AlertResource[T]) TypeName() string {
	// example: mezmo_threshold_alert
	return fmt.Sprintf("%s_%s_alert", PROVIDER_TYPE_NAME, r.typeName)
}

func (r *AlertResource[T]) NodeType() string {
	// example: threshold_alert
	return fmt.Sprintf("%s_alert", r.typeName)
}

func (r *AlertResource[T]) TerraformSchema() schema.Schema {
	return r.schema
}

func (r *AlertResource[T]) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("alert Value does not contain an interfaceable type")
	}

	alert, ok := component.Interface().(Alert)
	if !ok {
		return nil, errors.New("component Value cannot be cast to an Alert")
	}

	var model T
	r.toModelFunc(&model, &alert)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *AlertResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *AlertResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	stored, err := r.client.CreateAlert(r.getPipelineIdFunc(&plan).ValueString(), component, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating alert",
			"Could not create alert, unexpected error: "+err.Error(),
		)
		return
	}

	NullifyPlanFields(&plan, r.schema)

	r.toModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *AlertResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For this delete, we need the full alert struct to get field values that satisfy
	// the REST routes in `pipeline-service`. Normal deletes just use an id, but our api
	// paths for these are like `/v3/pipeline/:pipeline_id/:component_kind/:component_id/alert`
	alert, dd := r.fromModelFunc(&state, &state)
	if setDiagnosticsHasError(dd, &resp.Diagnostics) {
		return
	}

	err := r.client.DeleteAlert(r.getPipelineIdFunc(&state).ValueString(), alert, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Alert",
			"Could not delete alert, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *AlertResource[T]) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

// Read implements resource.Resource.
func (r *AlertResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Alert(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString(), ctx)
	// force re-creation of manually deleted resources
	if client.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Alert",
			fmt.Sprintf("Could not read alert with id %s and pipeline_id %s: %s",
				r.getIdFunc(&state), r.getPipelineIdFunc(&state), err.Error()),
		)
		return
	}

	NullifyPlanFields(&state, r.schema)

	r.toModelFunc(&state, component)
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *AlertResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	stored, err := r.client.UpdateAlert(r.getPipelineIdFunc(&state).ValueString(), component, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Alert",
			"Could not updated alert, unexpected error: "+err.Error(),
		)
		return
	}

	NullifyPlanFields(&plan, r.schema)

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *AlertResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema
}

func (r *AlertResource[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Malformed Resource Import Id",
			fmt.Sprintf(
				"The alert resource import id needs to be in the form of <pipeline id>/<alert id>. The "+
					"input \"%s\" did not contain two parts after parsing the id.",
				req.ID,
			),
		)
		return
	}

	pipelineId := strings.TrimSpace(parts[0])
	if pipelineId == "" {
		resp.Diagnostics.AddError(
			"Resource Import Id Missing Pipeline Id Part",
			"The pipeline id specified only contained whitespace.",
		)
	}

	alertId := strings.TrimSpace(parts[1])
	if alertId == "" {
		resp.Diagnostics.AddError(
			"Resource Import Id Missing Alert Id Part",
			"The alert id specified only contained whitespace",
		)
	}

	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), alertId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipeline_id"), pipelineId)...)
	}
}
