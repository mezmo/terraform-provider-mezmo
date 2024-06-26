package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	. "github.com/mezmo/terraform-provider-mezmo/internal/provider/models/modelutils"
	"golang.org/x/net/context"
)

type Client interface {
	Pipeline(id string, ctx context.Context) (*Pipeline, error)
	CreatePipeline(pipeline *Pipeline, ctx context.Context) (*Pipeline, error)
	UpdatePipeline(pipeline *Pipeline, ctx context.Context) (*Pipeline, error)
	DeletePipeline(id string, ctx context.Context) error

	Source(pipelineId string, id string, ctx context.Context) (*Source, error)
	CreateSource(pipelineId string, component *Source, ctx context.Context) (*Source, error)
	UpdateSource(pipelineId string, component *Source, ctx context.Context) (*Source, error)
	DeleteSource(pipelineId string, id string, ctx context.Context) error

	Destination(pipelineId string, id string, ctx context.Context) (*Destination, error)
	CreateDestination(pipelineId string, component *Destination, ctx context.Context) (*Destination, error)
	UpdateDestination(pipelineId string, component *Destination, ctx context.Context) (*Destination, error)
	DeleteDestination(pipelineId string, id string, ctx context.Context) error

	Processor(pipelineId string, id string, ctx context.Context) (*Processor, error)
	CreateProcessor(pipelineId string, component *Processor, ctx context.Context) (*Processor, error)
	UpdateProcessor(pipelineId string, component *Processor, ctx context.Context) (*Processor, error)
	DeleteProcessor(pipelineId string, id string, ctx context.Context) error

	Alert(pipelineId string, id string, ctx context.Context) (*Alert, error)
	CreateAlert(pipelineId string, alert *Alert, ctx context.Context) (*Alert, error) // POST
	UpdateAlert(pipelineId string, alert *Alert, ctx context.Context) (*Alert, error) // PUT
	DeleteAlert(pipelineId string, alert *Alert, ctx context.Context) error           // DELETE

	CreateAccessKey(accessKey *AccessKey, ctx context.Context) (*AccessKey, error)
	DeleteAccessKey(accessKey *AccessKey, ctx context.Context) error

	SharedSource(id string, ctx context.Context) (*SharedSource, error)
	CreateSharedSource(source *SharedSource, ctx context.Context) (*SharedSource, error)
	UpdateSharedSource(source *SharedSource, ctx context.Context) (*SharedSource, error)
	DeleteSharedSource(source *SharedSource, ctx context.Context) error

	PublishPipeline(pipelineId string, ctx context.Context) (*PublishPipeline, error)
}

func NewClient(endpoint string, authKey string, headers map[string]string) Client {
	return &client{
		httpClient: &http.Client{},
		endpoint:   endpoint,
		authKey:    authKey,
		headers:    headers,
	}
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

func readBody(result any, resp *http.Response, ctx context.Context) error {
	defer resp.Body.Close()
	bodyBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return newAPIError(resp.StatusCode, bodyBuffer, nil)
	}
	if result == nil || len(bodyBuffer) == 0 {
		return nil
	}
	// TF_LOG_PROVIDER=TRACE - Beware that this may print sensitive information
	err = json.Unmarshal(bodyBuffer, result)
	if err == nil {
		msg := Json(fmt.Sprintf("%s %s", resp.Request.Method, resp.Request.URL), result)
		tflog.Trace(ctx, msg)
	}
	return err
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
func (c *client) CreatePipeline(pipeline *Pipeline, ctx context.Context) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline", c.endpoint)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Pipeline]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	created := &envelope.Data
	return created, nil
}

// DeletePipeline implements Client.
func (c *client) DeletePipeline(id string, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// Pipeline implements Client.
func (c *client) Pipeline(id string, ctx context.Context) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Pipeline]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	pipeline := &envelope.Data
	return pipeline, nil
}

// UpdatePipeline implements Client.
func (c *client) UpdatePipeline(pipeline *Pipeline, ctx context.Context) (*Pipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s", c.endpoint, pipeline.Id)
	reqBody, err := json.Marshal(pipeline)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Pipeline]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	updated := &envelope.Data
	return updated, nil
}

// POST Source
func (c *client) CreateSource(pipelineId string, component *Source, ctx context.Context) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Source]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}

	source := &envelope.Data
	return source, nil
}

// DELETE Source
func (c *client) DeleteSource(pipelineId string, id string, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// GET Source
func (c *client) Source(pipelineId string, id string, ctx context.Context) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Source]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	source := &envelope.Data
	return source, nil
}

// PUT Source
func (c *client) UpdateSource(pipelineId string, component *Source, ctx context.Context) (*Source, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/source/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Source]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	source := &envelope.Data
	return source, nil
}

// GET Destination (sink)
func (c *client) Destination(pipelineId string, id string, ctx context.Context) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Destination]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	destination := &envelope.Data
	return destination, nil
}

// POST Destination (sink)
func (c *client) CreateDestination(pipelineId string, component *Destination, ctx context.Context) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Destination]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	destination := &envelope.Data
	return destination, nil
}

// DELETE Destination (sink)
func (c *client) DeleteDestination(pipelineId string, id string, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// PUT Destination (sink)
func (c *client) UpdateDestination(pipelineId string, component *Destination, ctx context.Context) (*Destination, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/sink/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Destination]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	destination := &envelope.Data
	return destination, nil
}

// GET Processor (transform)
func (c *client) Processor(pipelineId string, id string, ctx context.Context) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Processor]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}

// POST Processor (transform)
func (c *client) CreateProcessor(pipelineId string, component *Processor, ctx context.Context) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform", c.endpoint, pipelineId)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Processor]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}

// DELETE Processor (transform)
func (c *client) DeleteProcessor(pipelineId string, id string, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// PUT Processor (transform)
func (c *client) UpdateProcessor(pipelineId string, component *Processor, ctx context.Context) (*Processor, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/transform/%s", c.endpoint, pipelineId, component.Id)
	reqBody, err := json.Marshal(component)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Processor]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	processor := &envelope.Data
	return processor, nil
}

// GET Alert (not used by the UI)
func (c *client) Alert(pipelineId string, id string, ctx context.Context) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/alert/%s", c.endpoint, pipelineId, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Alert]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	alert := &envelope.Data
	return alert, nil
}

// POST Alert
func (c *client) CreateAlert(pipelineId string, alert *Alert, ctx context.Context) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId)
	reqBody, err := json.Marshal(alert)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Alert]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	createdAlert := &envelope.Data
	return createdAlert, nil
}

// PUT Alert
func (c *client) UpdateAlert(pipelineId string, alert *Alert, ctx context.Context) (*Alert, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert/%s", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId, alert.Id)
	reqBody, err := json.Marshal(alert)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[Alert]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	updatedAlert := &envelope.Data
	return updatedAlert, nil
}

// DELETE Alert
func (c *client) DeleteAlert(pipelineId string, alert *Alert, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/%s/%s/%s/alert/%s", c.endpoint, pipelineId, alert.ComponentKind, alert.ComponentId, alert.Id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// POST Access Key
func (c *client) CreateAccessKey(accessKey *AccessKey, ctx context.Context) (*AccessKey, error) {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s/access-key", c.endpoint, accessKey.SharedSourceId)
	reqBody, err := json.Marshal(accessKey)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[AccessKey]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	createdAccessKey := &envelope.Data
	return createdAccessKey, nil
}

// DELETE access key
func (c *client) DeleteAccessKey(accessKey *AccessKey, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s/access-key/%s", c.endpoint, accessKey.SharedSourceId, accessKey.Id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// POST Shared Source
func (c *client) CreateSharedSource(source *SharedSource, ctx context.Context) (*SharedSource, error) {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route", c.endpoint)
	reqBody, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[SharedSource]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	createdSharedSource := &envelope.Data
	return createdSharedSource, nil
}

// GET shared source
func (c *client) SharedSource(id string, ctx context.Context) (*SharedSource, error) {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s", c.endpoint, id)
	req := c.newRequest(http.MethodGet, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[SharedSource]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	source := &envelope.Data
	return source, nil
}

// PUT shared source
func (c *client) UpdateSharedSource(source *SharedSource, ctx context.Context) (*SharedSource, error) {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s", c.endpoint, source.Id)
	reqBody, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPut, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[SharedSource]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	updatedSharedSource := &envelope.Data
	return updatedSharedSource, nil
}

// DELETE shared source
func (c *client) DeleteSharedSource(source *SharedSource, ctx context.Context) error {
	url := fmt.Sprintf("%s/v3/pipeline/gateway-route/%s", c.endpoint, source.Id)
	req := c.newRequest(http.MethodDelete, url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	return readBody(nil, resp, ctx)
}

// POST publish pipeline
func (c *client) PublishPipeline(pipelineId string, ctx context.Context) (*PublishPipeline, error) {
	url := fmt.Sprintf("%s/v3/pipeline/%s/publish?allow_unconnected_edges=true", c.endpoint, pipelineId)
	// Because it's a POST, an empty body is required
	reqBody, err := json.Marshal(struct{}{})
	if err != nil {
		return nil, err
	}
	req := c.newRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	var envelope apiResponseEnvelope[PublishPipeline]
	if err := readBody(&envelope, resp, ctx); err != nil {
		return nil, err
	}
	created := &envelope.Data
	return created, nil
}
