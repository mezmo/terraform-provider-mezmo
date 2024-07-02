package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BlobFileConsolidationModel represents the configuration for consolidating
// files in a blob storage destination, e.g. Azure Blob Storage or S3.
type BlobFileConsolidationModel struct {
	Enabled             basetypes.BoolValue   `tfsdk:"enabled"`
	ProcessEverySeconds basetypes.Int64Value  `tfsdk:"process_every_seconds"`
	RequestSizeBytes    basetypes.Int64Value  `tfsdk:"requested_size_bytes"`
	BasePath            basetypes.StringValue `tfsdk:"base_path"`
}

// BlobFileConsolidationAttr is the schema for the BlobFileConsolidationModel.
var BlobFileConsolidationAttr = schema.SingleNestedAttribute{
	Optional:    true,
	Description: "This sink writes many small files out to azure blob storage. Enabling this process will allow the automatic consolidation of these small files into larger files of your choosing. This process will enable upon deployment and run on the chosen interval from `Processing Interval` creating files named `merged_[timestamp].log` where `timestamp` is the time since epoch when the actual file was created. The process will recursively access all files under the `Base Path`  to handle merging sub-directory logging structures.",
	Attributes: map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Optional:    true,
			Description: "Toggles whether the process is enabled.",
		},
		"process_every_seconds": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "How often to run the consolidation process in seconds",
			Default:     int64default.StaticInt64(600),
			Validators:  []validator.Int64{int64validator.Between(300, 3600)},
		},
		"requested_size_bytes": schema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "The requested size of the consolidated files in bytes.",
			Default:     int64default.StaticInt64(500_000_000),
			Validators:  []validator.Int64{int64validator.Between(50_000_000, 10_000_000_000)},
		},
		"base_path": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The path from the container to begin recursively looking for files to merge. A blank value indicates root. Merged files will only contain data from the respective folder.",
			Validators:  []validator.String{stringvalidator.LengthBetween(1, 512)},
		},
	},
}

// ToMap converts the BlobFileConsolidationModel to a map.
func (m *BlobFileConsolidationModel) ToMap() map[string]any {
	return map[string]any{
		"enabled":               m.Enabled.ValueBool(),
		"process_every_seconds": m.ProcessEverySeconds.ValueInt64(),
		"requested_size_bytes":  m.RequestSizeBytes.ValueInt64(),
		"base_path":             m.BasePath.ValueString(),
	}
}
