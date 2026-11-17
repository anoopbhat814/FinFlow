package domain

import "time"

type AuditEvent struct {
	ID         string            `json:"id"`
	ActorID    string            `json:"actorId,omitempty"`
	Action     string            `json:"action"`
	EntityType string            `json:"entityType"`
	EntityID   string            `json:"entityId"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	IP         string            `json:"ip,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
}
