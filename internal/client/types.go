package client

import (
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
