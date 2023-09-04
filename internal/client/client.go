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

	Destination(pipelineId string, id string) (*Destination, error)
	CreateDestination(pipelineId string, component *Destination) (*Destination, error)
	UpdateDestination(pipelineId string, component *Destination) (*Destination, error)
	DeleteDestination(pipelineId string, id string) error

	Processor(pipelineId string, id string) (*Processor, error)
	CreateProcessor(pipelineId string, component *Processor) (*Processor, error)
	UpdateProcessor(pipelineId string, component *Processor) (*Processor, error)
	DeleteProcessor(pipelineId string, id string) error
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

// Destination implements Client.
// Gets a destination.
func (c *client) Destination(pipelineId string, id string) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Destination]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	destination := &envelope.Data
	return destination, nil
}

// CreateDestination implements Client.
func (c *client) CreateDestination(pipelineId string, component *Destination) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Destination]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	destination := &envelope.Data
	return destination, nil
}

// DeleteDestination implements Client.
func (c *client) DeleteDestination(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// UpdateDestination implements Client.
func (c *client) UpdateDestination(pipelineId string, component *Destination) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Destination]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	destination := &envelope.Data
	return destination, nil
}

// Processor implements Client.
// Gets a Processor.
func (c *client) Processor(pipelineId string, id string) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Processor]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}

// CreateProcessor implements Client.
func (c *client) CreateProcessor(pipelineId string, component *Processor) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Processor]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}

// DeleteProcessor implements Client.
func (c *client) DeleteProcessor(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// UpdateProcessor implements Client.
func (c *client) UpdateProcessor(pipelineId string, component *Processor) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Processor]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}
