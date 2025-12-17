package destinations

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	. "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// BlobFileConsolidationAttr is the schema for the file consolidation model.
// This represents the configuration for consolidating files in a blob storage destination
// e.g. Azure Blob Storage or S3
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

// ToFileConsolidationObject converts a map representation of file consolidation configuration
// into a basetypes.ObjectValue suitable for use in the Terraform provider models.
func ToFileConsolidationObject(fc map[string]any) basetypes.ObjectValue {
	file_consolidation, _ := ObjectValue(map[string]attr.Type{
		"enabled":               BoolType,
		"process_every_seconds": Int64Type,
		"requested_size_bytes":  Int64Type,
		"base_path":             StringType,
	}, map[string]attr.Value{
		"enabled":               BoolValue(fc["enabled"].(bool)),
		"process_every_seconds": Int64Value(int64(fc["process_every_seconds"].(float64))),
		"requested_size_bytes":  Int64Value(int64(fc["requested_size_bytes"].(float64))),
		"base_path":             StringValue(fc["base_path"].(string)),
	})
	return file_consolidation
}
