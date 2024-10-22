package gohassapi

import (
	"time"
)

type apiCheck struct {
	Message string `json:"message"`
}

type ApiState struct {
	// entity_id, state, last_changed and attributes.
	EntityId string `json:"entity_id"`
	State string `json:"state"`
	LastChanged time.Time `json:"last_changed"`
	Attributes map[string]any `json:"attributes"`
}
