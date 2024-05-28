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

	Alert(pipelineId string, id string) (*Alert, error)
	CreateAlert(pipelineId string, alert *Alert) (*Alert, error) // POST
	UpdateAlert(pipelineId string, alert *Alert) (*Alert, error) // PUT
	DeleteAlert(pipelineId string, alert *Alert) error           // DELETE

	CreateAccessKey(accessKey *AccessKey) (*AccessKey, error)
	DeleteAccessKey(accessKey *AccessKey) error
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return newAPIError(resp)
	}

	return err
}

func readJson(result any, resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return newAPIError(resp)
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

// POST Source
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

// DELETE Source
func (c *client) DeleteSource(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// GET Source
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

// PUT Source
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

// GET Destination (sink)
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

// POST Destination (sink)
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

// DELETE Destination (sink)
func (c *client) DeleteDestination(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// PUT Destination (sink)
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

// GET Processor (transform)
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

// POST Processor (transform)
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

// DELETE Processor (transform)
func (c *client) DeleteProcessor(pipelineId string, id string) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// PUT Processor (transform)
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

// GET Alert (not used by the UI)
func (c *client) Alert(pipelineId string, id string) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/alert/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Alert]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	alert := &envelope.Data
	return alert, nil
}

// POST Alert
func (c *client) CreateAlert(pipelineId string, alert *Alert) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId)
	reqBody, err := json.Marshal(alert)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Alert]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	createdAlert := &envelope.Data
	return createdAlert, nil
}

// PUT Alert
func (c *client) UpdateAlert(pipelineId string, alert *Alert) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert/%s", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId, alert.Id)
	reqBody, err := json.Marshal(alert)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[Alert]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	updatedAlert := &envelope.Data
	return updatedAlert, nil
}

// DELETE Alert
func (c *client) DeleteAlert(pipelineId string, alert *Alert) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert/%s", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId, alert.Id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}

// POST Access Key
func (c *client) CreateAccessKey(accessKey *AccessKey) (*AccessKey, error) {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s/access-key", c.endpoint, accessKey.GatewayRouteId)
	reqBody, err := json.Marshal(accessKey)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	var envelope apiResponseEnvelope[AccessKey]
	if err := readJson(&envelope, resp, err); err != nil {
		return nil, err
	}
	createdAccessKey := &envelope.Data
	return createdAccessKey, nil
}

// DELETE access key
func (c *client) DeleteAccessKey(accessKey *AccessKey) error {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s/access-key/%s", c.endpoint, accessKey.GatewayRouteId, accessKey.Id)
	req := c.newRequest(http.MethodDelete, url, nil)
	return readBody(c.httpClient.Do(req))
}
