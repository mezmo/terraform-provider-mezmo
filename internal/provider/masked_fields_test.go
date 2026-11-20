package provider

import (
	"reflect"
	"testing"
)

func TestRestoreMaskedUserConfigFields(t *testing.T) {
	apiUserConfig := map[string]any{
		"api_key": "masked",
		"auth":    map[string]any{"password": "masked", "strategy": "basic"},
		"site":    "us1",
	}
	maskedFields := []string{"api_key", "auth.password"}
	userConfig := map[string]any{
		"api_key": "real-key",
		"auth":    map[string]string{"password": "real-password", "strategy": "basic"},
		"site":    "us1",
	}

	restoreMaskedUserConfigFields(apiUserConfig, maskedFields, userConfig)

	expected := map[string]any{
		"api_key": "real-key",
		"auth":    map[string]any{"password": "real-password", "strategy": "basic"},
		"site":    "us1",
	}
	if !reflect.DeepEqual(apiUserConfig, expected) {
		t.Fatalf("apiUserConfig = %v, want %v", apiUserConfig, expected)
	}
}

func TestRestoreMaskedUserConfigFieldsWithoutValue(t *testing.T) {
	apiUserConfig := map[string]any{
		"auth":   map[string]any{"password": "masked"},
		"secret": "masked",
		"token":  "masked",
	}
	maskedFields := []string{"auth.password", "secret", "token"}
	userConfig := map[string]any{
		"secret": nil,
	}

	restoreMaskedUserConfigFields(apiUserConfig, maskedFields, userConfig)

	expected := map[string]any{
		"auth":   map[string]any{"password": "masked"},
		"secret": "masked",
		"token":  "masked",
	}
	if !reflect.DeepEqual(apiUserConfig, expected) {
		t.Fatalf("apiUserConfig = %v, want %v", apiUserConfig, expected)
	}
}
