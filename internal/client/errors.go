package client

import (
	"fmt"
	"io"
	"net/http"
)

type apiResponseError struct {
	body   string
	status uint16
}

func newAPIError(resp *http.Response) apiResponseError {
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	status := uint16(resp.StatusCode)
	body := ""
	if err != nil {
		body = fmt.Sprintf("failed to read response body. reason: %s", err.Error())
	} else {
		body = string(respBody)
	}

	return apiResponseError{
		body:   body,
		status: status,
	}
}

// implement errors.Error interface
func (e apiResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.status, e.body)
}

func IsNotFoundError(target error) bool {
	err, ok := target.(apiResponseError)
	return ok && err.status == http.StatusNotFound
}
