package client

import (
	"time"

	. "github.com/mezmo-inc/terraform-provider-mezmo/internal/client/types"
)

type Client interface {
	Pipeline(id string) (*Pipeline, error)
	CreatePipeline(pipeline *Pipeline) (*Pipeline, error)
	UpdatePipeline(pipeline *Pipeline) (*Pipeline, error)
	DeletePipeline(id string) error
}

func NewClient() Client {
	return &noopClient{}
}

type noopClient struct {
}

// CreatePipeline implements Client.
func (*noopClient) CreatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	pipeline.UpdatedAt = time.Now()
	pipeline.Id = "Generated ID: " + time.Now().String()
	return pipeline, nil
}

// DeletePipeline implements Client.
func (*noopClient) DeletePipeline(id string) error {
	return nil
}

// Pipeline implements Client.
func (*noopClient) Pipeline(id string) (*Pipeline, error) {
	return &Pipeline{
		Id:        id,
		Title:     "Generated Title - " + id,
		UpdatedAt: time.Now(),
	}, nil
}

// UpdatePipeline implements Client.
func (*noopClient) UpdatePipeline(pipeline *Pipeline) (*Pipeline, error) {
	return pipeline, nil
}
