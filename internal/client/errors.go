package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type validationError struct {
	Path    string `json:"instancePath,omitempty"`
	Message string `json:"message,omitempty"`
}
type apiResponseError struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Status  uint16            `json:"status"`
	Errors  []validationError `json:"errors,omitempty"`
}

func newAPIError(resp *http.Response, _ context.Context) apiResponseError {
	defer resp.Body.Close()

	var result apiResponseError
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		result.Message = fmt.Sprintf("failed to read response body. reason: %s", err.Error())
		return result
	}
	return result
}

// implement errors.Error interface
func (e apiResponseError) Error() string {
	errString := fmt.Sprintf("%d %s: %s", e.Status, e.Code, e.Message)
	if len(e.Errors) > 0 {
		errors := make([]validationError, 0)
		for _, err := range e.Errors {
			if SkipThisError(err.Message) {
				continue
			}
			errors = append(errors, err)
		}
		jsonErrors, err := json.MarshalIndent(errors, "", "  ")
		if err == nil {
			errString += fmt.Sprintf("\nValidation Errors: %s", string(jsonErrors))
		}
	}
	return errString

}

func IsNotFoundError(target error) bool {
	err, ok := target.(apiResponseError)
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
