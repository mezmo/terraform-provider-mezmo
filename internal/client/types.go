package client

import (
	"time"
)

// Represents a Pipeline.
type Pipeline struct {
	Id        string     `json:"id,omitempty"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
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

// Represents an Alert, which is similar to a component, but not treated as such.
type Alert struct {
	Id            string         `json:"id,omitempty"`
	PipelineId    string         `json:"pipeline_id,omitempty"`
	ComponentKind string         `json:"component_kind,omitempty"`
	ComponentId   string         `json:"component_id,omitempty"`
	Inputs        []string       `json:"inputs,omitempty"`
	AlertConfig   map[string]any `json:"alert_config,omitempty"`
	Active        bool           `json:"active"`
}

type Source struct {
	BaseNode
	SharedSourceId string `json:"gateway_route_id,omitempty"`
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

type AccessKey struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	SharedSourceId string `json:"gateway_route_id"`
	Type           string `json:"type"`
	Key            string `json:"key,omitempty"`
}

type SharedSource struct {
	Id          string `json:"id"`
	ConsumerId  string `json:"consumer_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
}

type PublishPipeline struct {
	PipelineId string `json:"id"`
}
