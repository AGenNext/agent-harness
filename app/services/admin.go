package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminUser represents an admin user
type AdminUser struct {
	ID        string    `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`     // admin, developer, viewer
	Teams    []string `json:"teams"`    // teams allowed
	Skills   []string `json:"skills"`   // skills allowed
	Created  time.Time `json:"created"`
	LastLogin time.Time `json:"last_login"`
}

// AdminService manages admin users
type AdminService struct {
	users   map[string]*AdminUser
	secret  string
}

// NewAdminService creates a new admin service
func NewAdminService() *AdminService {
	return &AdminService{
		users:  make(map[string]*AdminUser),
		secret: os.Getenv("HARNESS_SECRET"),
	}
}

// InitDefaultAdmin creates default admin user
func (s *AdminService) InitDefaultAdmin() {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASS")
	adminUser := os.Getenv("ADMIN_USER")

	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	if adminUser == "" {
		adminUser = "admin"
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)

	s.users[adminEmail] = &AdminUser{
		ID:       "admin-1",
		Email:    adminEmail,
		Name:     adminUser,
		Role:     "admin",
		Teams:    []string{"*"},
		Skills:   []string{"*"},
		Created:  time.Now(),
	}

	fmt.Printf("✅ Admin user created: %s\n", adminEmail)
}

// Authenticate verifies admin credentials
func (s *AdminService) Authenticate(ctx context.Context, email, password string) (*AdminUser, error) {
	user, ok := s.users[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	user.LastLogin = time.Now()
	return user, nil
}

// Authorize checks if user has permission for action
func (s *AdminService) Authorize(user *AdminUser, agent, skill string) bool {
	if user.Role == "admin" {
		return true
	}

	for _, t := range user.Teams {
		if t == "*" || t == agent {
			for _, s := range user.Skills {
				if s == "*" || s == skill {
					return true
				}
			}
		}
	}
	return false
}

// RequireAuth middleware protects routes
func (s *AdminService) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify token (simplified - use proper JWT in production)
		if auth != s.secret {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}