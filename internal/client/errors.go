package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func newAPIError(status int, bodyBuffer []byte, _ context.Context) ApiResponseError {
	result := ApiResponseError{
		Status: uint16(status),
	}
	if len(bodyBuffer) == 0 {
		return result
	}
	err := json.Unmarshal(bodyBuffer, &result)
	if err != nil {
		result.Message = fmt.Sprintf("failed to decode JSON body. reason: %s", err.Error())
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
	default:
		return false
	}
}
