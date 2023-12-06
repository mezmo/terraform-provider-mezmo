package client

import (
	"encoding/json"
	"slices"
	"time"
)

var DISALLOWED_EMPTY_STRING_FIELDS_BY_PARSER = map[string][]string{
	"parse_key_value": {"field_delimiter", "key_delimiter"},
}

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

func filterEmptyStrings(m map[string]any) map[string]any {
	parserName, ok := m["parser"].(string)
	if !ok {
		return m
	}
	filteredFields, ok := DISALLOWED_EMPTY_STRING_FIELDS_BY_PARSER[parserName]
	if !ok {
		return m
	}
	userOptions, _ := m["options"].(map[string]any)
	for key, val := range userOptions {
		if slices.Contains(filteredFields, key) {
			switch v := val.(type) {
			case string:
				if v == "" {
					delete(userOptions, key)
				}
			}
		}
	}
	if userOptions != nil {
		m["options"] = userOptions
	}
	return m
}

func (u *BaseNode) MarshalJSON() ([]byte, error) {
	u.UserConfig = filterEmptyStrings(u.UserConfig)

	parsers, _ := u.UserConfig["parsers"].([]map[string]any)
	for _, parserAttrs := range parsers {
		filterEmptyStrings(parserAttrs)
	}
	if parsers != nil {
		u.UserConfig["parsers"] = parsers
	}

	type Alias BaseNode
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	})
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
