package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mezmo/terraform-provider-mezmo/v4/internal/provider/models/modelutils"
)

type validationError struct {
	Path    string `json:"instancePath,omitempty"`
	Message string `json:"message,omitempty"`
}
type ApiResponseError struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Status  uint16            `json:"status"`
	Errors  []validationError `json:"errors,omitempty"`
}

func newAPIError(ctx context.Context, status int, bodyBuffer []byte, _ context.Context) ApiResponseError {
	result := ApiResponseError{
		Status: uint16(status),
	}
	if len(bodyBuffer) == 0 {
		return result
	}
	// Errors *should* follow the expected format, but to be safe, we'll unmarshal into an
	// interface and assert the fields we expect.
	var e map[string]any

	if err := json.Unmarshal(bodyBuffer, &e); err != nil {
		result.Code = "EUNKNOWN"
		result.Message = "There was an error, but the response from the server was not understood."
		result.Errors = append(result.Errors, validationError{
			Path:    "raw error body",
			Message: string(bodyBuffer),
		})
		return result
	}

	msg := modelutils.Json("API Error response", e)
	tflog.Trace(ctx, msg)

	if code, ok := e["code"].(string); ok {
		result.Code = code
	}
	if message, ok := e["message"].(string); ok {
		result.Message = message
	}
	if errors, ok := e["errors"].([]any); ok {
		for _, err := range errors {
			if e, ok := err.(map[string]any); ok {
				vError := validationError{}
				if path, ok := e["instancePath"].(string); ok {
					vError.Path = path
				}
				if message, ok := e["message"].(string); ok {
					vError.Message = message
				}
				result.Errors = append(result.Errors, vError)
			}
		}
	} else if errorString, ok := e["errors"].(string); ok {
		result.Errors = append(result.Errors, validationError{
			Message: errorString,
		})
	}
	return result
}

// implement errors.Error interface
func (e ApiResponseError) Error() string {
	errString := fmt.Sprintf("%d %s: %s", e.Status, e.Code, e.Message)
	if len(e.Errors) > 0 {
		errors := make([]validationError, 0)
		for _, err := range e.Errors {
			if SkipThisError(err.Message) {
				continue
			}
			errors = append(errors, err)
		}
		jsonErrors, err := JSONMarshal(errors)
		if err == nil {
			errString += fmt.Sprintf("\nValidation Errors: %s", string(jsonErrors))
		}
	}
	return errString
}

// We cannot use `json.MarshalIndent` because it will escape HTML by default.
// An Encoder will let us control that behavior, and add the indenting we want.
func JSONMarshal(t any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func IsNotFoundError(target error) bool {
	err, ok := target.(ApiResponseError)
	return ok && err.Status == http.StatusNotFound
}

// Validation errors can be wordy and redundant. Skip ones that don't make sense
// ourside of the schema's context.
func SkipThisError(msg string) bool {
	switch msg {
	case "must match \"then\" schema":
		return true
	case "must match \"else\" schema":
		return true
	case "must match a schema in anyOf":
		return true
	case "must match a schema in allOf":
		return true
	case "must match a schema in oneOf":
		return true
	default:
		return false
	}
}
