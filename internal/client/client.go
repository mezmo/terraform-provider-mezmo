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

	Source(pipelineId string, id string) (*Source, error)
	CreateSource(pipelineId string, component *Source) (*Source, error)
	UpdateSource(pipelineId string, component *Source) (*Source, error)
	DeleteSource(pipelineId string, id string) error

	Sink(pipelineId string, id string) (*Sink, error)
	CreateSink(pipelineId string, component *Sink) (*Sink, error)
	UpdateSink(pipelineId string, component *Sink) (*Sink, error)
	DeleteSink(pipelineId string, id string) error

	Transform(pipelineId string, id string) (*Transform, error)
	CreateTransform(pipelineId string, component *Transform) (*Transform, error)
	UpdateTransform(pipelineId string, component *Transform) (*Transform, error)
	DeleteTransform(pipelineId string, id string) error
}

func NewClient(endpoint string, authKey string, headers map[string]string) Client {
	return &client{
		httpClient: &http.Client{},
		endpoint:   endpoint,
		authKey:    authKey,
		headers:    headers,
	}
}

// Envelope is {meta, data}, but meta is not used
type apiResponseEnvelope[T any] struct {
	Data T `json:"data"`
}

type client struct {
	httpClient *http.Client
	endpoint   string
	authKey    string
	headers    map[string]string
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
	if c.authKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", c.authKey))
	}
	if len(c.headers) > 0 {
		for k, v := range c.headers {
			req.Header.Add(k, v)
		}
	}

	if method != http.MethodGet {
		req.Header.Add("Content-Type", "application/json")
	}
	return req
}

// CreateSource implements Client.
func (c *client) CreateSource(pipelineId string, component *Source) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Source]
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
func (c *client) Source(pipelineId string, id string) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Source]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	source := &envelope.Data
	return source, nil
}

// UpdateSource implements Client.
func (c *client) UpdateSource(pipelineId string, component *Source) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Source]
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
func (c *client) Sink(pipelineId string, id string) (*Sink, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Sink]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	sink := &envelope.Data
	return sink, nil
}

// CreateSink implements Client.
func (c *client) CreateSink(pipelineId string, component *Sink) (*Sink, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Sink]
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
func (c *client) UpdateSink(pipelineId string, component *Sink) (*Sink, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Sink]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	sink := &envelope.Data
	return sink, nil
}

// Transform implements Client.
// Gets a Transform.
func (c *client) Transform(pipelineId string, id string) (*Transform, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Transform]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	transform := &envelope.Data
	return transform, nil
}

// CreateTransform implements Client.
func (c *client) CreateTransform(pipelineId string, component *Transform) (*Transform, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Transform]
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
func (c *client) UpdateTransform(pipelineId string, component *Transform) (*Transform, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Transform]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	transform := &envelope.Data
	return transform, nil
}
