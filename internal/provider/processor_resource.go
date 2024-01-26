package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/processors"
)

type ProcessorModel interface {
	AggregateV2ProcessorModel |
		CompactFieldsProcessorModel |
		DecryptFieldsProcessorModel |
		DedupeProcessorModel |
		DropFieldsProcessorModel |
		EncryptFieldsProcessorModel |
		EventToMetricProcessorModel |
		FilterProcessorModel |
		FlattenFieldsProcessorModel |
		MapFieldsProcessorModel |
		MetricsTagCardinalityLimitProcessorModel |
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
	client   Client
	typeName string // The name to use as part of the terraform resource name: mezmo_{typeName}_processor
	// the corresponding pipeline service node type
	nodeName          string
	fromModelFunc     processorFromModelFunc[T]
	toModelFunc       processorToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	schema            schema.Schema
}

func (r *ProcessorResource[T]) TypeName() string {
	return fmt.Sprintf("%s_%s_processor", PROVIDER_TYPE_NAME, r.typeName)
}

func (r *ProcessorResource[T]) NodeType() string {
	// conversion to underscore and appending suffix
	// disambiguates resources with similar node types
	// example: http sink and http source
	return strings.ReplaceAll(r.nodeName, "-", "_") + "_processor"
}

func (r *ProcessorResource[T]) TerraformSchema() schema.Schema {
	return r.schema
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

	NullifyPlanFields(&plan, r.schema)

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
func (r *ProcessorResource[T]) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
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

	NullifyPlanFields(&state, r.schema)

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

	NullifyPlanFields(&plan, r.schema)

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *ProcessorResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema
}

func (r *ProcessorResource[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Malformed Resource Import Id",
			fmt.Sprintf(
				"The processor resource import id needs to be in the form of <pipeline id>/<processor id>. The "+
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

	nodeId := strings.TrimSpace(parts[1])
	if nodeId == "" {
		resp.Diagnostics.AddError(
			"Resource Import Id Missing Processor Id Part",
			"The processor id specified only contained whitespace",
		)
	}

	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), nodeId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipeline_id"), pipelineId)...)
	}
}
