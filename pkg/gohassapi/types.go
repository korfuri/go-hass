package gohassapi

import (
	"time"
)

type Check struct {
	Message string `json:"message"`
}

type State struct {
	// entity_id, state, last_changed and attributes.
	EntityId string `json:"entity_id"`
	State string `json:"state"`
	LastChanged time.Time `json:"last_changed"`
	Attributes map[string]any `json:"attributes"`
}

type Service struct {
	Name string `json:"name"`
	Description string `json:"description"`
	// TODO: the service schema is very different in practice from
	// what is documented, and appears very flexible. Figure out
	// what the actual schema is before describing it here.
}

type ServiceDomain struct {
	Domain string `json:"domain"`
	Services map[string]Service `json:"services"`
}
