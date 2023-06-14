package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client/types"
)

type Client interface {
	Pipeline(id string) (*Pipeline, error)
	CreatePipeline(pipeline *Pipeline) (*Pipeline, error)
	UpdatePipeline(pipeline *Pipeline) (*Pipeline, error)
	DeletePipeline(id string) error
}

func NewClient(endpoint string, authKey string, authHeader string, authAdditional string) Client {
	return &client{
		httpClient:     &http.Client{},
		authKey:        authKey,
		authHeader:     authHeader,
		endpoint:       endpoint,
		authAdditional: authAdditional,
	}
}

type client struct {
	httpClient     *http.Client
	authKey        string
	authHeader     string
	endpoint       string
	authAdditional string
}

// CreatePipeline implements Client.
func (c *client) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v1/pipelines", c.endpoint)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var stored Pipeline
	if err := readJson(&stored, resp, err); err != nil {
		return nil, err
	}
	return &stored, nil
}

// DeletePipeline implements Client.
func (c *client) DeletePipeline(id string) error {
	url := fmt.Sprintf("%s/v1/pipelines/%s", c.endpoint, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	_, err := readBody(c.httpClient.Do(req))
	return err
}

// Pipeline implements Client.
func (c *client) Pipeline(id string) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v1/pipelines/%s", c.endpoint, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var stored Pipeline
	if err := readJson(&stored, resp, err); err != nil {
		return nil, err
	}
	return &stored, nil
}

// UpdatePipeline implements Client.
func (c *client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v1/pipelines/%s", c.endpoint, pipeline.Id)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var stored Pipeline
	if err := readJson(&stored, resp, err); err != nil {
		return nil, err
	}
	return &stored, nil
}

func readBody(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	return body, err
}

func readJson(result any, resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *client) newRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		// Creating a request only fails if the method is invalid
		panic(err)
	}
	req.Header.Add(c.authHeader, c.authKey)
	if c.authHeader == "x-auth-account-id" {
		req.Header.Add("x-auth-user-email", c.authAdditional)
	}

	if method != http.MethodGet {
		req.Header.Add("Content-Type", "application/json")
	}
	return req
}
