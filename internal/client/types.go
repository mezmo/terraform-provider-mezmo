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
type Component struct {
	Id           string         `json:"id,omitempty"`
	Type         string         `json:"type"`
	Inputs       []string       `json:"inputs,omitempty"`
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
	UserConfig   map[string]any `json:"user_config"`
	GenerationId int64          `json:"generation_id"`
}

// Represents a full pipeline response from the service.
type pipelineResponse struct {
	Id         string      `json:"id"`
	Sources    []Component `json:"sources"`
	Transforms []Component `json:"transforms"`
	Sinks      []Component `json:"sinks"`
}

func (p *pipelineResponse) findSource(id string) (*Component, error) {
	for _, s := range p.Sources {
		if s.Id == id {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("Source %s not found in pipeline %s", id, p.Id)
}
