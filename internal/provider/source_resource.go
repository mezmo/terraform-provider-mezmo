package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/sources"
)

type SourceModel interface {
	AgentSourceModel |
		AzureEventHubSourceModel |
		DemoSourceModel |
		FluentSourceModel |
		HttpSourceModel |
		KafkaSourceModel |
		KinesisFirehoseSourceModel |
		LogAnalysisSourceModel |
		LogStashSourceModel |
		OpenTelemetryTracesSourceModel |
		PrometheusRemoteWriteSourceModel |
		S3SourceModel |
		SplunkHecSourceModel |
		SQSSourceModel
}

type SourceResource[T SourceModel] struct {
	client            Client
	typeName          string // The name to use as part of the terraform resource name: mezmo_{typeName}_source
	fromModelFunc     sourceFromModelFunc[T]
	toModelFunc       sourceToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	getSchemaFunc     getSchemaFunc
}

func (r *SourceResource[T]) TypeName() string {
	return r.typeName
}

func (r *SourceResource[T]) TerraformSchema() schema.Schema {
	return r.getSchemaFunc()
}

func (r *SourceResource[T]) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("component Value does not contain an interfaceable type")
	}

	source, ok := component.Interface().(Source)
	if !ok {
		return nil, errors.New("component Value cannot be cast to Source")
	}

	var model T
	r.toModelFunc(&model, &source)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *SourceResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *SourceResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	if os.Getenv("DEBUG_SOURCE") == "1" {
		PrintJSON("----- Destination TO Create api ---", component)
	}

	stored, err := r.client.CreateSource(r.getPipelineIdFunc(&plan).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating source",
			"Could not create source, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_SOURCE") == "1" {
		PrintJSON("----- Destination FROM Create api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *SourceResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteSource(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not source, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *SourceResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName + "_source"
}

// Read implements resource.Resource.
func (r *SourceResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Source(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			fmt.Sprintf("Could not read source with id %s and pipeline_id %s: %s",
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
func (r *SourceResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	if os.Getenv("DEBUG_SOURCE") == "1" {
		PrintJSON("----- Destination TO Update api ---", component)
	}

	stored, err := r.client.UpdateSource(r.getPipelineIdFunc(&state).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			"Could not updated source, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_SOURCE") == "1" {
		PrintJSON("----- Destination FROM Update api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *SourceResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.getSchemaFunc()
}
