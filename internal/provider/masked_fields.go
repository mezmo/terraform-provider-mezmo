package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// restoreMaskedUserConfigFields puts the real values back into a user_config the API
// returned with its secrets masked, so state keeps what the user configured.
//
//   - apiUserConfig is the user_config from the API response. It is updated in place.
//   - maskedFields lists the dotted paths the API masked, from the response's
//     masked_fields.
//   - userConfig holds the real values, built from the plan or from prior state.
//
// A masked field with no value in userConfig keeps the masked value.
func restoreMaskedUserConfigFields(apiUserConfig map[string]any, maskedFields []string, userConfig map[string]any) {
	if len(maskedFields) == 0 {
		return
	}
	// Models build some nested values as map[string]string; a JSON round trip gives
	// the map[string]any shape the API response has.
	raw, err := json.Marshal(userConfig)
	if err != nil {
		return
	}
	var values map[string]any
	if err := json.Unmarshal(raw, &values); err != nil {
		return
	}
	for _, field := range maskedFields {
		value, ok := userConfigValue(values, field)
		if !ok || value == nil {
			continue
		}
		setUserConfigValue(apiUserConfig, field, value)
	}
}

// isImportedState reports whether state holds only the ids set by an import, which
// leaves no user_config to restore masked values from.
func isImportedState(ctx context.Context, state tfsdk.State) bool {
	var generationId types.Int64
	state.GetAttribute(ctx, path.Root("generation_id"), &generationId)
	return generationId.IsNull()
}

func userConfigValue(config map[string]any, field string) (any, bool) {
	var value any = config
	for _, key := range strings.Split(field, ".") {
		object, ok := value.(map[string]any)
		if !ok {
			return nil, false
		}
		if value, ok = object[key]; !ok {
			return nil, false
		}
	}
	return value, true
}

func setUserConfigValue(config map[string]any, field string, value any) {
	keys := strings.Split(field, ".")
	object := config
	for _, key := range keys[:len(keys)-1] {
		next, ok := object[key].(map[string]any)
		if !ok {
			return
		}
		object = next
	}
	object[keys[len(keys)-1]] = value
}
