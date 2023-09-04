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

type Processor struct {
	BaseNode
	Outputs []struct {
		Id    string `json:"id"`
		Label string `json:"label"`
	} `json:"outputs,omitempty"`
}

type Destination struct {
	BaseNode
}

// Represents a full pipeline response from the service.
type pipelineResponse struct {
	Id           string        `json:"id"`
	Sources      []Source      `json:"sources"`
	Processors   []Processor   `json:"transforms"`
	Destinations []Destination `json:"sinks"`
}
