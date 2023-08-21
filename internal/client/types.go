package client

import (
	"fmt"
	"time"
)

// Represents a Pipeline.
type Pipeline struct {
	Id        string     `json:"id,omitempty"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// Represents a source, processor or destination.
type BaseNode struct {
	Id           string         `json:"id,omitempty"`
	Type         string         `json:"type"`
	Inputs       []string       `json:"inputs,omitempty"`
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
	UserConfig   map[string]any `json:"user_config"`
	GenerationId int64          `json:"generation_id"`
}

type Source struct {
	BaseNode
	GatewayRouteId string `json:"gateway_route_id,omitempty"`
}
type Transform struct {
	BaseNode
	OutputNames []struct {
		Id    string `json:"id"`
		Label string `json:"label"`
	} `json:"output_names,omitempty"`
}
type Sink struct {
	BaseNode
}

// Represents a full pipeline response from the service.
type pipelineResponse struct {
	Id         string      `json:"id"`
	Sources    []Source    `json:"sources"`
	Transforms []Transform `json:"transforms"`
	Sinks      []Sink      `json:"sinks"`
}

func (p *pipelineResponse) findSource(id string) (*Source, error) {
	for _, s := range p.Sources {
		if s.Id == id {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Source %s not found in pipeline %s", id, p.Id)
}

func (p *pipelineResponse) findSink(id string) (*Sink, error) {
	for _, s := range p.Sinks {
		if s.Id == id {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Sink %s not found in pipeline %s", id, p.Id)
}

func (p *pipelineResponse) findTransform(id string) (*Transform, error) {
	for _, s := range p.Transforms {
		if s.Id == id {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Transform %s not found in pipeline %s", id, p.Id)
}
