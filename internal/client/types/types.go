package types

import "time"

// Represents a Pipeline.
type Pipeline struct {
	Id        string     `json:"id,omitempty"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// Represents a source, processor or destination.
type Component struct {
	Id          string         `json:"id,omitempty"`
	Type        string         `json:"type"`
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	UserOptions map[string]any `json:"user_options"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
}
