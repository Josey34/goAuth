package middleware_test

import (
	"net/http/httptest"
	"testing"

	"goauth/middleware"

	"github.com/gin-gonic/gin"
)

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		requiredRole string
		expectError  bool
	}{
		{
			name:         "user has required role",
			userRole:     "admin",
			requiredRole: "admin",
			expectError:  false,
		},
		{
			name:         "user lacks required role",
			userRole:     "user",
			requiredRole: "admin",
			expectError:  true,
		},
		{
			name:         "user role exact match",
			userRole:     "user",
			requiredRole: "user",
			expectError:  false,
		},
		{
			name:         "empty user role",
			userRole:     "",
			requiredRole: "admin",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleMiddleware := middleware.RequireRole(tt.requiredRole)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin", nil)

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Set the user role in context (simulating Auth middleware)
			c.Set("userRole", tt.userRole)

			// Call middleware
			roleMiddleware(c)

			// Check if context was aborted for error cases
			if tt.expectError && !c.IsAborted() {
				t.Errorf("expected context to be aborted for error case")
			}
			if !tt.expectError && c.IsAborted() {
				t.Errorf("expected context NOT to be aborted for success case")
			}

			// For error cases, verify response was written
			if tt.expectError && w.Body.Len() == 0 {
				t.Errorf("expected error response but got empty body")
			}
		})
	}
}

func TestRequireRoleMultipleRoles(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		requiredRole string
		shouldPass   bool
	}{
		{
			name:         "admin can access admin endpoint",
			userRole:     "admin",
			requiredRole: "admin",
			shouldPass:   true,
		},
		{
			name:         "user cannot access admin endpoint",
			userRole:     "user",
			requiredRole: "admin",
			shouldPass:   false,
		},
		{
			name:         "moderator cannot access admin endpoint",
			userRole:     "moderator",
			requiredRole: "admin",
			shouldPass:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleMiddleware := middleware.RequireRole(tt.requiredRole)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admin/dashboard", nil)

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userRole", tt.userRole)

			roleMiddleware(c)

			if tt.shouldPass && c.IsAborted() {
				t.Errorf("expected access to be granted but was denied")
			}
			if !tt.shouldPass && !c.IsAborted() {
				t.Errorf("expected access to be denied but was granted")
			}
		})
	}
}
