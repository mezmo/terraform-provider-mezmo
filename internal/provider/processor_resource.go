package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/processors"
)

type ProcessorModel interface {
	CompactFieldsProcessorModel |
		DecryptFieldsProcessorModel |
		DedupeProcessorModel |
		DropFieldsProcessorModel |
		EncryptFieldsProcessorModel |
		FlattenFieldsProcessorModel |
		MapFieldsProcessorModel |
		ParseProcessorModel |
		ParseSequentiallyProcessorModel |
		ReduceProcessorModel |
		RouteProcessorModel |
		SampleProcessorModel |
		ScriptExecutionProcessorModel |
		StringifyProcessorModel |
		UnrollProcessorModel
}

type ProcessorResource[T ProcessorModel] struct {
	client            Client
	typeName          string // The name to use as part of the terraform resource name: mezmo_{typeName}_processor
	fromModelFunc     processorFromModelFunc[T]
	toModelFunc       processorToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	getSchemaFunc     getSchemaFunc
}

func (r *ProcessorResource[T]) TypeName() string {
	return r.typeName + "_processor"
}

func (r *ProcessorResource[T]) TerraformSchema() schema.Schema {
	return r.getSchemaFunc()
}

func (r *ProcessorResource[T]) ConvertToTerraformModel(component *reflect.Value) (*reflect.Value, error) {
	if !component.CanInterface() {
		return nil, errors.New("component Value does not contain an interfaceable type")
	}

	processor, ok := component.Interface().(Processor)
	if !ok {
		return nil, errors.New("component Value cannot be cast to Processor")
	}

	var model T
	r.toModelFunc(&model, &processor)
	res := reflect.ValueOf(model)
	return &res, nil
}

// Configure implements resource.ResourceWithConfigure.
func (r *ProcessorResource[T]) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Processor Configure Type",
			fmt.Sprintf("Expected client.Client, got: %T. Please report this issue to Mezmo.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *ProcessorResource[T]) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	if os.Getenv("DEBUG_PROCESSOR") == "1" {
		PrintJSON("----- Processor TO Create api ---", component)
	}

	stored, err := r.client.CreateProcessor(r.getPipelineIdFunc(&plan).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating processor",
			"Could not create processor, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_PROCESSOR") == "1" {
		PrintJSON("----- Processor FROM Create api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *ProcessorResource[T]) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state T
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteProcessor(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting processor",
			"Could not delete processor, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *ProcessorResource[T]) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TypeName()
}

// Read implements resource.Resource.
func (r *ProcessorResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Processor(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading processor",
			fmt.Sprintf("Could not read processor with id %s and pipeline_id %s: %s",
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
func (r *ProcessorResource[T]) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	if os.Getenv("DEBUG_PROCESSOR") == "1" {
		PrintJSON("----- Processor TO Update api ---", component)
	}

	stored, err := r.client.UpdateProcessor(r.getPipelineIdFunc(&state).ValueString(), component)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating processor",
			"Could not update processor, unexpected error: "+err.Error(),
		)
		return
	}

	if os.Getenv("DEBUG_PROCESSOR") == "1" {
		PrintJSON("----- Processor FROM Update api ---", stored)
	}

	NullifyPlanFields(&plan, r.getSchemaFunc())

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *ProcessorResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.getSchemaFunc()
}
