package service

import (
	"context"
	"fmt"
)

// =============================================
// Multi-Tenancy for Agent Harness
// =============================================

// Tenant represents an organization/team
type Tenant struct {
	ID           string            `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"` // URL-friendly
	OwnerID     int64           `json:"owner_id"`
	Settings    TenantSettings   `json:"settings"`
	CreatedAt  int64           `json:"created_at"`
	UpdatedAt  int64           `json:"updated_at"`
}

// TenantSettings for tenant-specific config
type TenantSettings struct {
	// Agents
	AllowedAgents []string `json:"allowed_agents"` // Which agents this tenant can use
	AgentQuota   int      `json:"agent_quota"`   // Max concurrent agents

	// Resources
	MaxBuilds   int `json:"max_builds"`   // Per day
	MaxStorage  int `json:"max_storage"`  // MB

	// Notifications
	NotifySlack    bool `json:"notify_slack"`
	NotifyDiscord  bool `json:"notify_discord"`
	NotifyMattermost bool `json:"notify_mattermost"`

	// Billing
	Plan    string `json:"plan"` // free, pro, enterprise
	BillingEmail string `json:"billing_email"`
}

// DefaultTenantSettings returns default settings
func DefaultTenantSettings() TenantSettings {
	return TenantSettings{
		AllowedAgents: []string{"code-assist", "code-review", "code-tester", "code-deploy"},
		AgentQuota:    5,
		MaxBuilds:     100,
		MaxStorage:    1000, // 1GB
		Plan:         "free",
	}
}

// =============================================
// Tenant Manager
// =============================================

type TenantManager struct {
	tenants map[string]*Tenant
}

func NewTenantManager() *TenantManager {
	return &TenantManager{
		tenants: make(map[string]*Tenant),
	}
}

// CreateTenant creates a new tenant
func (m *TenantManager) CreateTenant(ctx context.Context, name, slug, ownerEmail string) (*Tenant, error) {
	// Check slug unique
	if _, exists := m.tenants[slug]; exists {
		return nil, fmt.Errorf("tenant slug already exists: %s", slug)
	}

	tenant := &Tenant{
		ID:        fmt.Sprintf("tenant_%d", len(m.tenants)+1),
		Name:      name,
		Slug:     slug,
		OwnerID:  1, // Would come from auth
		Settings: DefaultTenantSettings(),
	}

	m.tenants[slug] = tenant
	return tenant, nil
}

// GetTenant retrieves a tenant by slug
func (m *TenantManager) GetTenant(slug string) (*Tenant, error) {
	t, ok := m.tenants[slug]
	if !ok {
		return nil, fmt.Errorf("tenant not found: %s", slug)
	}
	return t, nil
}

// UpdateTenant updates tenant settings
func (m *TenantManager) UpdateTenant(slug string, settings TenantSettings) error {
	if t, ok := m.tenants[slug]; ok {
		t.Settings = settings
		return nil
	}
	return fmt.Errorf("tenant not found: %s", slug)
}

// ListTenants lists all tenants
func (m *TenantManager) ListTenants() []*Tenant {
	list := make([]*Tenant, 0, len(m.tenants))
	for _, t := range m.tenants {
		list = append(list, t)
	}
	return list
}

// =============================================
// Tenant Isolation
// =============================================

// GetTenantID extracts tenant from context
func GetTenantID(ctx context.Context) string {
	// Would get from JWT/context
	return "default"
}

// CanAccessAgent checks if tenant can use agent
func (t *Tenant) CanAccessAgent(agent string) bool {
	for _, a := range t.Settings.AllowedAgents {
		if a == agent {
			return true
		}
	}
	return false
}

// CheckQuota checks if under agent quota
func (t *Tenant) CheckQuota(currentAgents int) bool {
	return currentAgents < t.Settings.AgentQuota
}

// =============================================
// Context with Tenant
// =============================================

type tenantKey string

const (
	ctxTenant tenantKey = "tenant"
)

// WithTenant adds tenant to context
func WithTenant(ctx context.Context, tenant *Tenant) context.Context {
	return context.WithValue(ctx, ctxTenant, tenant)
}

// FromContext gets tenant from context
func FromContext(ctx context.Context) (*Tenant, error) {
	t, ok := ctx.Value(ctxTenant).(*Tenant)
	if !ok {
		return nil, fmt.Errorf("no tenant in context")
	}
	return t, nil
}

// =============================================
// Usage Example
// =============================================

/*
// API Usage:

// Create tenant
POST /api/v1/tenants
{
  "name": "Acme Corp",
  "slug": "acme",
  "plan": "pro"
}

// List tenant's agents (isolated)
GET /api/v1/agents?tenant=acme

// Tenant-specific settings  
PUT /api/v1/tenants/acme/settings
{
  "allowed_agents": ["code-assist", "code-review"],
  "agent_quota": 2,
  "plan": "enterprise"
}

// Enforce in API:
agent := ctrl.GetAgent(r.Context(), name)
tenant, _ := service.FromContext(r.Context())
if !tenant.CanAccessAgent(agent.Name) {
    return errors.New("agent not allowed for tenant")
}
*/