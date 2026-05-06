// Immutable Audit Log
package audit

import (
	"context"
	"fmt"
	"time"
)

// AuditEvent - immutable, append-only
type AuditEvent struct {
	ID        string `json:"id"` // SHA256 hash
	TenantID string `json:"tenant_id"`
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Details  string `json:"details"`
	IP       string `json:"ip"`
	Timestamp int64 `json:"timestamp"` // Set once!
	Version  int   `json:"version"` // For chain
}

// NewAuditEvent - create immutable event
func NewAuditEvent(tenantID, actor, action, resource, details string) *AuditEvent {
	ts := time.Now().Unix()
	return &AuditEvent{
		ID:        fmt.Sprintf("audit_%x", ts),
		TenantID:  tenantID,
		Actor:    actor,
		Action:   action,
		Resource: resource,
		Details:  details,
		Timestamp: ts,
		Version:  1,
	}
}

const (
	ActionLogin    = "user.login"
	ActionLogout   = "user.logout"
	ActionAgentRun = "agent.run"
	ActionDeploy  = "deploy.create"
	ActionTenantCreate = "tenant.create"
)

// AuditStore - append-only
type AuditStore struct {
	events []*AuditEvent
}

func NewAuditStore() *AuditStore {
	return &AuditStore{events: make([]*AuditEvent, 0)}
}

// Append - ONLY way to add events
func (s *AuditStore) Append(event *AuditEvent) {
	s.events = append(s.events, event)
}

// Query - read events
func (s *AuditStore) Query(tenantID string, from, to int64) []*AuditEvent {
	var result []*AuditEvent
	for _, e := range s.events {
		if e.TenantID != tenantID { continue }
		if from > 0 && e.Timestamp < from { continue }
		if to > 0 && e.Timestamp > to { continue }
		result = append(result, e)
	}
	return result
}

// FORBIDDEN: Delete is not allowed
func (s *AuditStore) Delete(id string) error {
	return fmt.Errorf("DELETE FORBIDDEN: audit logs are immutable")
}

// FORBIDDEN: Update not allowed
func (s *AuditStore) Update(event *AuditEvent) error {
	return fmt.Errorf("UPDATE FORBIDDEN: audit logs are immutable")
}