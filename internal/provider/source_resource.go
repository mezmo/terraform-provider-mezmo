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
	. "github.com/mezmo/terraform-provider-mezmo/v5/internal/provider/models/sources"
)

type SourceModel interface {
	AgentSourceModel |
		AzureEventHubSourceModel |
		DatadogSourceModel |
		DemoSourceModel |
		FluentSourceModel |
		HttpSourceModel |
		KafkaSourceModel |
		KinesisFirehoseSourceModel |
		LogAnalysisSourceModel |
		LogAnalysisIngestionSourceModel |
		LogStashSourceModel |
		OpenTelemetryLogsSourceModel |
		OpenTelemetryMetricsSourceModel |
		OpenTelemetryTracesSourceModel |
		PrometheusRemoteWriteSourceModel |
		S3SourceModel |
		SplunkHecSourceModel |
		SQSSourceModel |
		WebhookSourceModel
}

type SourceResource[T SourceModel] struct {
	client            Client
	typeName          string // The name to use as part of the terraform resource name: mezmo_{typeName}_source
	nodeName          string // the corresponding pipeline service node type
	fromModelFunc     sourceFromModelFunc[T]
	toModelFunc       sourceToModelFunc[T]
	getIdFunc         idGetterFunc[T]
	getPipelineIdFunc idGetterFunc[T]
	schema            schema.Schema
}

func (r *SourceResource[T]) TypeName() string {
	return fmt.Sprintf("%s_%s_source", PROVIDER_TYPE_NAME, r.typeName)
}

func (r *SourceResource[T]) NodeType() string {
	// conversion to underscore and appending suffix
	// disambiguates resources with similar node types
	// example: http sink and http source
	return strings.ReplaceAll(r.nodeName, "-", "_") + "_source"
}

func (r *SourceResource[T]) TerraformSchema() schema.Schema {
	return r.schema
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

	stored, err := r.client.CreateSource(r.getPipelineIdFunc(&plan).ValueString(), component, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating source",
			"Could not create source, unexpected error: "+err.Error(),
		)
		return
	}

	NullifyPlanFields(&plan, r.schema)

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
	err := r.client.DeleteSource(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString(), ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source",
			"Could not delete source, unexpected error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *SourceResource[T]) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.TypeName()
}

// Read implements resource.Resource.
func (r *SourceResource[T]) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state T
	if diags := req.State.Get(ctx, &state); setDiagnosticsHasError(diags, &resp.Diagnostics) {
		return
	}

	component, err := r.client.Source(r.getPipelineIdFunc(&state).ValueString(), r.getIdFunc(&state).ValueString(), ctx)
	// force re-creation of manually deleted resources
	if client.IsNotFoundError(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Source",
			fmt.Sprintf("Could not read source with id %s and pipeline_id %s: %s",
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

	stored, err := r.client.UpdateSource(r.getPipelineIdFunc(&state).ValueString(), component, ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source",
			"Could not update source, unexpected error: "+err.Error(),
		)
		return
	}

	NullifyPlanFields(&plan, r.schema)

	r.toModelFunc(&plan, stored)
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema implements resource.Resource.
func (r *SourceResource[T]) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.schema
}

func (r *SourceResource[T]) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Malformed Resource Import Id",
			fmt.Sprintf(
				"The source resource import id needs to be in the form of <pipeline id>/<source id>. The "+
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
			"Resource Import Id Missing Source Id Part",
			"The source id specified only contained whitespace",
		)
	}

	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), nodeId)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipeline_id"), pipelineId)...)
	}
}
