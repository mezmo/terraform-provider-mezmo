package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/destinations"
	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/provider/models/modelutils"
)

type DestinationModel interface {
	AzureBlobStorageDestinationModel |
		BlackholeDestinationModel |
		DatadogLogsDestinationModel |
		DatadogMetricsDestinationModel |
		ElasticSearchDestinationModel |
		GcpCloudStorageDestinationModel |
		HoneycombLogsDestinationModel |
		HttpDestinationModel |
		KafkaDestinationModel |
		LokiDestinationModel |
		MezmoDestinationModel |
		NewRelicDestinationModel |
		PrometheusRemoteWriteDestinationModel |
		S3DestinationModel |
		SplunkHecLogsDestinationModel
}

type DestinationResource[T DestinationModel] struct {
	client            Client
	typeName          string // The name to use as part of the terraform resource name: mezmo_{typeName}_destination
	fromModelFunc     destinationFromModelFunc[T]
	toModelFunc       destinationToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	getSchemaFunc     getSchemaFunc
}

func (r *DestinationResource[T]) TypeName() string {
	return r.typeName
}

func (r *DestinationResource[T]) TerraformSchema() schema.Schema {
	return r.getSchemaFunc()
}

func (r *DestinationResource[T]) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("component Value does not contain an interfaceable type")
	}

	dest, ok := component.Interface().(Destination)
	if !ok {
		return nil, errors.New("component Value cannot be cast to Destination")
	}

	var model T
	r.toModelFunc(&model, &dest)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *DestinationResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Destination Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to Mezmo.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *DestinationResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	if os.Getenv("DEBUG_DESTINATION") == "1" {
		PrintJSON("----- Destination TO Create api ---", component)
	}

	stored, err := r.client.CreateDestination(r.getPipelineIdFunc(&plan).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating destination",
			"Could not create destination, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_DESTINATION") == "1" {
		PrintJSON("----- Destination FROM Create api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *DestinationResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDestination(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting destination",
			"Could not delete destination, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *DestinationResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName + "_destination"
}

// Read implements resource.Resource.
func (r *DestinationResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Destination(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading destination",
			fmt.Sprintf("Could not read destination with id %s and pipeline_id %s: %s",
				r.getIdFunc(&state), r.getPipelineIdFunc(&state), err.Error()),
		)
		return
	}

	NullifyPlanFields(&state, r.getSchemaFunc())

	r.toModelFunc(&state, component)
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *DestinationResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	if os.Getenv("DEBUG_DESTINATION") == "1" {
		PrintJSON("----- Destination TO Update api ---", component)
	}

	// Set id from the current state (not in plan)
	stored, err := r.client.UpdateDestination(r.getPipelineIdFunc(&state).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating destination",
			"Could not updated destination, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_DESTINATION") == "1" {
		PrintJSON("----- Destination FROM Update api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *DestinationResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.getSchemaFunc()
}
