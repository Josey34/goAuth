package middleware_test

import (
	"net/http/httptest"
	"testing"

	"goauth/middleware"
	mock_test "goauth/tests/mock"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		shouldFail  bool
		expectError bool
	}{
		{
			name:        "valid token",
			authHeader:  "Bearer valid-token",
			shouldFail:  false,
			expectError: false,
		},
		{
			name:        "missing auth header",
			authHeader:  "",
			shouldFail:  true,
			expectError: true,
		},
		{
			name:        "invalid bearer format",
			authHeader:  "InvalidFormat token",
			shouldFail:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockToken := mock_test.NewMockTokenService()
			if tt.shouldFail && tt.name == "invalid token" {
				mockToken.ValidateShouldFail = true
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authMiddleware := middleware.Auth(mockToken)
			authMiddleware(c)

			if tt.expectError && !c.IsAborted() {
				t.Errorf("expected context to be aborted for error case")
			}

			if !tt.expectError && c.GetString("userID") == "" {
				t.Errorf("expected userID in context but got empty")
			}

			if tt.expectError && w.Body.Len() == 0 {
				t.Errorf("expected error response but got empty body")
			}
		})
	}
}
