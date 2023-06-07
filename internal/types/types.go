package types

import "time"

type Pipeline struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}
