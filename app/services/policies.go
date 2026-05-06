package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

// ============== Policy Types ==============

// Policy defines access rules
type Policy struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Effect     string   `json:"effect"`   // allow, deny
	Priority   int      `json:"priority"`  // Higher = first
	Subjects   []string `json:"subjects"`  // user, team, role
	Resources  []string `json:"resources"` // agents, workflows
	Actions    []string `json:"actions"`   // read, write, execute
	Conditions map[string]string `json:"conditions,omitempty"`
}

// RateLimitPolicy rate limiting
type RateLimitPolicy struct {
	Window     time.Duration `json:"window"`     // 1m, 5m, 1h
	MaxRequests int        `json:"max_requests"`
	Burst     int         `json:"burst"`
}

// SecurityPolicy security rules
type SecurityPolicy struct {
	RequireApproval  bool     `json:"require_approval"`
	AllowedEnvs     []string `json:"allowed_envs"`
	AllowedSkills  []string `json:"allowed_skills"`
	BlockPatterns  []string `json:"block_patterns"` // Blocked commands
	AutoTimeout   int      `json:"auto_timeout"` // seconds
}

// ============== Policy Store ==============

type PolicyStore struct {
	policies       map[string]*Policy
	rateLimits    map[string]*RateLimitPolicy
	security     *SecurityPolicy
}

// NewPolicyStore creates policy store
func NewPolicyStore() *PolicyStore {
	return &PolicyStore{
		policies:    make(map[string]*Policy),
		rateLimits: make(map[string]*RateLimitPolicy),
		security:  DefaultSecurityPolicy(),
	}
}

// DefaultPolicies returns default policies
func DefaultPolicies() map[string]*Policy {
	return map[string]*Policy{
		"admin-all": {
			Name: "admin-all", Effect: "allow", Priority: 100,
			Subjects: []string{"role:admin"},
			Resources: []string{"*"},
			Actions: []string{"*"},
		},
		"dev-code-assist": {
			Name: "dev-code-assist", Effect: "allow", Priority: 50,
			Subjects: []string{"role:developer"},
			Resources: []string{"code-assist"},
			Actions: []string{"read", "execute"},
		},
		"dev-deploy": {
			Name: "dev-deploy", Effect: "allow", Priority: 40,
			Subjects: []string{"role:developer"},
			Resources: []string{"code-deploy"},
			Actions: []string{"read"},
			Conditions: map[string]string{"env": "dev,staging"},
		},
		"viewer-read": {
			Name: "viewer-read", Effect: "allow", Priority: 10,
			Subjects: []string{"role:viewer"},
			Resources: []string{"*"},
			Actions: []string{"read"},
		},
		"default-deny": {
			Name: "default-deny", Effect: "deny", Priority: 0,
			Subjects: []string{"*"},
			Resources: []string{"*"},
			Actions: []string{"*"},
		},
	}
}

// DefaultSecurityPolicy returns security defaults
func DefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		RequireApproval:  true,
		AllowedEnvs:    []string{"dev", "staging", "prod"},
		AllowedSkills:  []string{"fix", "review", "test", "deploy"},
		BlockPatterns: []string{
			"rm -rf /",
			"DROP TABLE",
			"curl | sh",
			"exec(*",
		},
		AutoTimeout: 300,
	}
}

// ============== Policy Methods ==============

func (s *PolicyStore) AddPolicy(p *Policy) {
	s.policies[p.Name] = p
}

func (s *PolicyStore) GetPolicy(name string) *Policy {
	return s.policies[name]
}

func (s *PolicyStore) ListPolicies() []*Policy {
	list := make([]*Policy, 0, len(s.policies))
	for _, p := range s.policies {
		list = append(list, p)
	}
	return list
}

// CheckAccess evaluates policies
func (s *PolicyStore) CheckAccess(subject, resource, action string, conditions map[string]string) bool {
	var bestPolicy *Policy

	for _, p := range s.policies {
		if !matchesSubjects(p.Subjects, subject) {
			continue
		}
		if !matchesResources(p.Resources, resource) {
			continue
		}
		if !matchesActions(p.Actions, action) {
			continue
		}

		// Check conditions
		pass := true
		for k, v := range p.Conditions {
			if conditions[k] != v {
				pass = false
				break
			}
		}

		if pass && (bestPolicy == nil || p.Priority > bestPolicy.Priority) {
			bestPolicy = p
		}
	}

	if bestPolicy == nil {
		return false
	}
	return bestPolicy.Effect == "allow"
}

// ============== Matching ==============

func matchesSubjects(subjects []string, subject string) bool {
	for _, s := range subjects {
		if s == "*" || s == subject {
			return true
		}
		if hasPrefix(s, "role:") && hasPrefix(subject, "role:") {
			return true
		}
	}
	return false
}

func matchesResources(resources []string, resource string) bool {
	for _, r := range resources {
		if r == "*" || r == resource {
			return true
		}
	}
	return false
}

func matchesActions(actions []string, action string) bool {
	for _, a := range actions {
		if a == "*" || a == action {
			return true
		}
	}
	return false
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// ============== Rate Limiting ==============

func (s *PolicyStore) SetRateLimit(key string, policy *RateLimitPolicy) {
	s.rateLimits[key] = policy
}

func (s *PolicyStore) CheckRateLimit(key string) (allowed bool, remaining int) {
	policy, ok := s.rateLimits[key]
	if !ok {
		return true, 999
	}

	// Simplified - use Redis in production
	now := time.Now()
	elapsed := now.Sub(now) // Placeholder

	// Check burst first
	if policy.Burst > 0 {
		return true, policy.Burst
	}

	if policy.MaxRequests > 0 {
		return true, policy.MaxRequests
	}

	return true, policy.MaxRequests
}

// ============== Security ==============

func (s *PolicyStore) Security() *SecurityPolicy {
	return s.security
}

func (s *SecurityPolicy) ValidateCommand(cmd string) error {
	for _, pattern := range s.BlockPatterns {
		if contains(cmd, pattern) {
			return fmt.Errorf("command blocked: %s matches %s", cmd, pattern)
		}
	}
	return nil
}

func (s *SecurityPolicy) ValidateEnv(env string) error {
	for _, e := range s.AllowedEnvs {
		if e == env {
			return nil
		}
	}
	return fmt.Errorf("environment not allowed: %s", env)
}

func (s *SecurityPolicy) ValidateSkill(skill string) error {
	for _, s := range s.AllowedSkills {
		if s == skill {
			return nil
		}
	}
	return fmt.Errorf("skill not allowed: %s", skill)
}

// ============== Middleware ==============

// PolicyMiddleware enforces policies
func PolicyMiddleware(store *PolicyStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get subject from token
			subject := r.Header.Get("X-Subject")
			resource := r.URL.Path
			action := r.Method

			if !store.CheckAccess(subject, resource, action, nil) {
				http.Error(w, "forbidden by policy", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}