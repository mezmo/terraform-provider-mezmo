package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/mezmo/terraform-provider-mezmo/internal/client"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
)

const AZURE_BLOB_STORAGE_DESTINATION_TYPE_NAME = "azure_blob_storage"
const AZURE_BLOB_STORAGE_DESTINATION_NODE_NAME = "azure-blob-storage"

type AzureBlobStorageDestinationModel struct {
	Id                  String                      `tfsdk:"id"`
	PipelineId          String                      `tfsdk:"pipeline_id"`
	Title               String                      `tfsdk:"title"`
	Description         String                      `tfsdk:"description"`
	Inputs              List                        `tfsdk:"inputs"`
	GenerationId        Int64                       `tfsdk:"generation_id"`
	AckEnabled          Bool                        `tfsdk:"ack_enabled" user_config:"true"`
	BatchTimeoutSeconds Int64                       `tfsdk:"batch_timeout_secs" user_config:"true"`
	Encoding            String                      `tfsdk:"encoding" user_config:"true"`
	Compression         String                      `tfsdk:"compression" user_config:"true"`
	ContainerName       String                      `tfsdk:"container_name" user_config:"true"`
	ConnectionString    String                      `tfsdk:"connection_string" user_config:"true"`
	Prefix              String                      `tfsdk:"prefix" user_config:"true"`
	FileConsolidation   *BlobFileConsolidationModel `tfsdk:"file_consolidation"`
}

var AzureBlobStorageResourceSchema = schema.Schema{
	Description: "Publishes events to Azure Blob Storage",
	Attributes: ExtendBaseAttributes(map[string]schema.Attribute{
		"encoding": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The encoding to apply to the data",
			Default:     stringdefault.StaticString("text"),
			Validators:  []validator.String{stringvalidator.OneOf("json", "text")},
		},
		"compression": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The compression strategy used on the encoded data prior to sending",
			Default:     stringdefault.StaticString("none"),
			Validators:  []validator.String{stringvalidator.OneOf("gzip", "none")},
		},
		"container_name": schema.StringAttribute{
			Required:    true,
			Description: "The name of the container for blob storage",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"connection_string": schema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "A connection string for the account that contains an access key",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"prefix": schema.StringAttribute{
			Optional:    true,
			Description: "A prefix to be applied to all object keys",
			Validators:  []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"file_consolidation": BlobFileConsolidationAttr,
	}, []string{"batch_timeout_secs"}),
}

func AzureBlobStorageFromModel(plan *AzureBlobStorageDestinationModel, previousState *AzureBlobStorageDestinationModel) (*Destination, diag.Diagnostics) {
	dd := diag.Diagnostics{}

	component := Destination{
		BaseNode: BaseNode{
			Type:        AZURE_BLOB_STORAGE_DESTINATION_NODE_NAME,
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Inputs:      StringListValueToStringSlice(plan.Inputs),
			UserConfig: map[string]any{
				"ack_enabled":        plan.AckEnabled.ValueBool(),
				"batch_timeout_secs": plan.BatchTimeoutSeconds.ValueInt64(),
				"encoding":           plan.Encoding.ValueString(),
				"compression":        plan.Compression.ValueString(),
				"container_name":     plan.ContainerName.ValueString(),
				"connection_string":  plan.ConnectionString.ValueString(),
			},
		},
	}

	if !plan.Prefix.IsNull() {
		component.UserConfig["prefix"] = plan.Prefix.ValueString()
	}

	if plan.FileConsolidation != nil {
		component.UserConfig["file_consolidation"] = plan.FileConsolidation.ToMap()
	}

	if previousState != nil {
		component.Id = previousState.Id.ValueString()
		component.GenerationId = previousState.GenerationId.ValueInt64()
	}

	return &component, dd
}

func AzureBlobStorageToModel(plan *AzureBlobStorageDestinationModel, component *Destination) {
	plan.Id = StringValue(component.Id)
	if component.Title != "" {
		plan.Title = StringValue(component.Title)
	}
	if component.Description != "" {
		plan.Description = StringValue(component.Description)
	}
	plan.GenerationId = Int64Value(component.GenerationId)
	plan.Inputs = SliceToStringListValue(component.Inputs)
	plan.AckEnabled = BoolValue(component.UserConfig["ack_enabled"].(bool))
	plan.BatchTimeoutSeconds = Int64Value(int64(component.UserConfig["batch_timeout_secs"].(float64)))
	plan.Encoding = StringValue(component.UserConfig["encoding"].(string))
	plan.Compression = StringValue(component.UserConfig["compression"].(string))
	plan.ContainerName = StringValue(component.UserConfig["container_name"].(string))
	plan.ConnectionString = StringValue(component.UserConfig["connection_string"].(string))
	if component.UserConfig["prefix"] != nil {
		plan.Prefix = StringValue(component.UserConfig["prefix"].(string))
	}

	if fc, ok := component.UserConfig["file_consolidation"].(map[string]any); ok {
		plan.FileConsolidation = &BlobFileConsolidationModel{
			Enabled:             BoolValue(fc["enabled"].(bool)),
			ProcessEverySeconds: Int64Value(int64(fc["process_every_seconds"].(float64))),
			RequestSizeBytes:    Int64Value(int64(fc["requested_size_bytes"].(float64))),
			BasePath:            StringValue(fc["base_path"].(string)),
		}
	}
}
