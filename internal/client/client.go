package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	Pipeline(id string) (*Pipeline, error)
	CreatePipeline(pipeline *Pipeline) (*Pipeline, error)
	UpdatePipeline(pipeline *Pipeline) (*Pipeline, error)
	DeletePipeline(id string) error

	Source(pipelineId string, id string) (*Component, error)
	CreateSource(pipelineId string, component *Component) (*Component, error)
	UpdateSource(pipelineId string, component *Component) (*Component, error)
	DeleteSource(pipelineId string, id string) error

	Sink(pipelineId string, id string) (*Component, error)
	CreateSink(pipelineId string, component *Component) (*Component, error)
	UpdateSink(pipelineId string, component *Component) (*Component, error)
	DeleteSink(pipelineId string, id string) error

	Transform(pipelineId string, id string) (*Component, error)
	CreateTransform(pipelineId string, component *Component) (*Component, error)
	UpdateTransform(pipelineId string, component *Component) (*Component, error)
	DeleteTransform(pipelineId string, id string) error
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

// Envelope is {meta, data}, but meta is not used
type apiResponseEnvelope[T any] struct {
	Data T `json:"data"`
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
	url := fmt.Sprintf("%s/v3/pipeline", c.endpoint)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Pipeline]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	created := &envelope.Data
	return created, nil
}

// DeletePipeline implements Client.
func (c *client) DeletePipeline(id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// Pipeline implements Client.
func (c *client) Pipeline(id string) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Pipeline]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	pipeline := &envelope.Data
	return pipeline, nil
}

// UpdatePipeline implements Client.
func (c *client) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, pipeline.Id)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Pipeline]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	updated := &envelope.Data
	return updated, nil
}

func readBody(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	return err
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

// CreateSource implements Client.
func (c *client) CreateSource(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	source := &envelope.Data
	return source, nil
}

// DeleteSource implements Client.
func (c *client) DeleteSource(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// Source implements Client.
// Gets a source from a pipeline.
func (c *client) Source(pipelineId string, id string) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, pipelineId)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[pipelineResponse]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	pipeline := &envelope.Data
	return pipeline.findSource(id)
}

// UpdateSource implements Client.
func (c *client) UpdateSource(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	source := &envelope.Data
	return source, nil
}

// TODO: The following methods need an implementation
// using Source(), CreateSource(), UpdateSource(), DeleteSource() as examples.

// Sink implements Client.
// Gets a sink.
func (c *client) Sink(pipelineId string, id string) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, pipelineId)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[pipelineResponse]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	pipeline := &envelope.Data
	return pipeline.findSink(id)
}

// CreateSink implements Client.
func (c *client) CreateSink(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	sink := &envelope.Data
	return sink, nil
}

// DeleteSink implements Client.
func (c *client) DeleteSink(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// UpdateSink implements Client.
func (c *client) UpdateSink(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	sink := &envelope.Data
	return sink, nil
}

// Transform implements Client.
// Gets a Transform.
func (c *client) Transform(pipelineId string, id string) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, pipelineId)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[pipelineResponse]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	pipeline := &envelope.Data
	return pipeline.findTransform(id)
}

// CreateTransform implements Client.
func (c *client) CreateTransform(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	transform := &envelope.Data
	return transform, nil
}

// DeleteTransform implements Client.
func (c *client) DeleteTransform(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// UpdateTransform implements Client.
func (c *client) UpdateTransform(pipelineId string, component *Component) (*Component, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Component]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	transform := &envelope.Data
	return transform, nil
}
