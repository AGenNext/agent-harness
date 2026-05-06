package services

import (
	"context"
	"fmt"
	"strings"
)

// ============== Fine-Grained Access ==============

// SubjectType - user, team, service, role
type SubjectType string

const (
	SubjectUser    SubjectType = "user"
	SubjectTeam  SubjectType = "team"
	SubjectRole SubjectType = "role"
	SubjectService SubjectType = "service"
)

// ResourceType - repo, agent, workflow, etc
type ResourceType string

const (
	ResourceRepo    ResourceType = "repo"
	ResourceAgent  ResourceType = "agent"
	ResourceWorkflow ResourceType = "workflow"
	ResourceSecret ResourceType = "secret"
	ResourceEnv   ResourceType = "env"
)

// Action - what can be done
type Action string

const (
	ActionRead    Action = "read"
	ActionWrite  Action = "write"
	ActionExecute Action = "execute"
	ActionDelete Action = "delete"
	ActionAdmin  Action = "admin"
)

// FGAEntry - fine-grained access entry
type FGAEntry struct {
	Subject    string       `json:"subject"`     // user:john, team:eng, role:admin
	Relation  string       `json:"relation"`   // owner, editor, viewer, executor
	Resource  string       `json:"resource"`  // repo:*, agent:*
	ResourceType ResourceType `json:"resource_type"`
	Actions   []Action     `json:"actions"`   // read, write, execute
	Conditions string     `json:"conditions,omitempty"`
}

// FGAStore - Fine-Grained Access store
type FGAStore struct {
	entries   map[string][]*FGAEntry
	relations map[string][]string
}

// NewFGAStore creates FGA store
func NewFGAStore() *FGAStore {
	s := &FGAStore{
		entries:   make(map[string][]*FGAEntry),
		relations: make(map[string][]string),
	}
	s.initDefaults()
	return s
}

// Default FGA entries
func (s *FGAStore) initDefaults() {
	// User permissions
	s.Grant("user:admin", "owner", "repo:*", ResourceRepo, []Action{ActionRead, ActionWrite, ActionExecute, ActionDelete, ActionAdmin})
	s.Grant("user:admin", "admin", "agent:*", ResourceAgent, []Action{ActionRead, ActionWrite, ActionExecute, ActionAdmin})
	s.Grant("user:admin", "admin", "workflow:*", ResourceWorkflow, []Action{ActionRead, ActionWrite, ActionExecute, ActionAdmin})

	// Team permissions
	s.Grant("team:eng", "editor", "repo:*", ResourceRepo, []Action{ActionRead, ActionWrite, ActionExecute})
	s.Grant("team:eng", "viewer", "repo:public/*", ResourceRepo, []Action{ActionRead})

	// Role permissions
	s.Grant("role:developer", "executor", "agent:code-assist", ResourceAgent, []Action{ActionRead, ActionExecute})
	s.Grant("role:developer", "executor", "agent:code-review", ResourceAgent, []Action{ActionRead, ActionExecute})
	s.Grant("role:developer", "executor", "agent:code-tester", ResourceAgent, []Action{ActionRead, ActionExecute})
	s.Grant("role:developer", "viewer", "workflow:*", ResourceWorkflow, []Action{ActionRead, ActionExecute})

	s.Grant("role:viewer", "viewer", "repo:*", ResourceRepo, []Action{ActionRead})
	s.Grant("role:viewer", "viewer", "agent:*", ResourceAgent, []Action{ActionRead})
}

// Grant permission
func (s *FGAStore) Grant(subject, relation, resource string, resourceType ResourceType, actions []Action) {
	key := fmt.Sprintf("%s:%s", subject, resource)
	s.entries[key] = append(s.entries[key], &FGAEntry{
		Subject:      subject,
		Relation:    relation,
		Resource:   resource,
		ResourceType: resourceType,
		Actions:    actions,
	})
}

// Check permission
func (s *FGAStore) Check(subject, relation, resource string, action Action) bool {
	key := fmt.Sprintf("%s:%s", subject, resource)
	entries, ok := s.entries[key]
	if !ok {
		// Try wildcard
		wildcard := fmt.Sprintf("%s:*", strings.Split(subject, ":")[0])
		key := fmt.Sprintf("%s:*", resource)
		entries, ok = s.entries[key]
		if !ok {
			key = fmt.Sprintf("%s:*", wildcard)
			entries, ok = s.entries[key]
		}
	}

	if !ok {
		return false
	}

	for _, e := range entries {
		if matchesRelation(e.Relation, relation) {
			for _, a := range e.Actions {
				if a == Action || a == ActionAdmin {
					return true
				}
			}
		}
	}
	return false
}

// CheckRelation checks subject has relation to resource
func (s *FGAStore) CheckRelation(subject, relation, resource string) bool {
	key := fmt.Sprintf("%s:%s", subject, resource)
	if entries, ok := s.entries[key]; ok {
		for _, e := range entries {
			if e.Relation == relation {
				return true
			}
		}
	}
	return false
}

// GetRelations returns all relations for a resource
func (s *FGAStore) GetRelations(resource string) []string {
	var rels []string
	for key := range s.entries {
		if strings.HasPrefix(key, ":"+resource) || strings.HasSuffix(key, ":"+resource) {
			if entries, ok := s.entries[key]; ok {
				for _, e := range entries {
					rels = append(rels, e.Relation)
				}
			}
		}
	}
	return rels
}

// ListPermissions lists all permissions for subject
func (s *FGAStore) ListPermissions(subject string) []*FGAEntry {
	var list []*FGAEntry
	for key, entries := range s.entries {
		if strings.HasPrefix(key, subject+":") {
			list = append(list, entries...)
		}
	}
	return list
}

// Revoke removes permission
func (s *FGAStore) Revoke(subject, relation, resource string) {
	key := fmt.Sprintf("%s:%s", subject, resource)
	delete(s.entries, key)
}

// ============== Helpers ==============

func matchesRelation(granted, requested string) bool {
	if granted == requested {
		return true
	}
	// Hierarchy: owner > editor > executor > viewer
	 hierarchy := map[string]int{
		"owner":   4,
		"editor": 3,
		"executor": 2,
		"viewer": 1,
	}
	g, ok := hierarchy[granted]
	if !ok {
		return false
	}
	r, ok := hierarchy[requested]
	if !ok {
		return false
	}
	return g >= r
}

// ============== FGA Model ==============

// Object - who owns what
type Object struct {
	Type  ResourceType `json:"type"`
	ID    string       `json:"id"`
	Owner string       `json:"owner"` // user:john
}

// GetOwner returns owner of object
func (o *Object) GetOwner() string {
	return o.Owner
}

// CanAccess checks if subject can access object
func (o *Object) CanAccess(store *FGAStore, subject string, action Action) bool {
	// Direct owner
	if store.Check(fmt.Sprintf("user:%s", subject), "owner", fmt.Sprintf("%s:%s", o.Type, o.ID), action) {
		return true
	}
	// Relation check
	return store.Check(subject, "viewer", fmt.Sprintf("%s:%s", o.Type, o.ID), action)
}

// ============== Tuples (for debug) ==============

// ListTuples returns all tuples
func (s *FGAStore) ListTuples() []string {
	var tuples []string
	for key, entries := range s.entries {
		for _, e := range entries {
			tuples = append(tuples, fmt.Sprintf("%s:%s -> %s", e.Subject, e.Relation, key))
		}
	}
	return tuples
}