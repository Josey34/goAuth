package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goauth/config"
	"goauth/dto"
	"goauth/factory"
	"goauth/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func setupTestServer(t *testing.T) (*httptest.Server, *factory.Factory, func()) {
	cfg := &config.Config{
		Port:           8080,
		DBPath:         ":memory:",
		JWTSecret:      "test-secret-key-for-jwt-tokens-12345",
		AccessTTL:      15 * time.Minute,
		RefreshTTL:     7 * 24 * time.Hour,
		BcryptCost:     10,
		AllowedOrigins: []string{"*"},
		LogLevel:       "info",
		RateLimitRPS:   100,
		RateLimitBurst: 200,
	}

	f, err := factory.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create factory: %v", err)
	}

	if err := initializeSchema(f.DB); err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	zerolog.SetGlobalLevel(zerolog.Disabled)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.Logger(zlog.Logger))
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	r.Use(rateLimiter.Limit())

	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", f.AuthHandler.Register)
	auth.POST("/login", f.AuthHandler.Login)
	auth.POST("/refresh", f.AuthHandler.Refresh)

	protected := api.Group("/auth")
	protected.Use(middleware.Auth(f.TokenService))
	protected.GET("/profile", f.UserHandler.GetProfile)
	protected.PUT("/profile", f.UserHandler.UpdateProfile)

	admin := api.Group("/admin")
	admin.Use(middleware.Auth(f.TokenService))
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/dashboard", f.UserHandler.AdminDashboard)

	server := httptest.NewServer(r)

	cleanup := func() {
		server.Close()
		f.DB.Close()
	}

	return server, f, cleanup
}

func initializeSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	return err
}

func makeRequest(t *testing.T, method, path string, body interface{}, token *string, baseURL string) (*http.Response, map[string]interface{}) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", baseURL, path), bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil && *token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	}
	resp.Body.Close()

	return resp, result
}

func TestCompleteAuthFlow(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("register new user", func(t *testing.T) {
		registerReq := dto.RegisterRequest{
			Email:    "alice@example.com",
			Password: "SecurePassword123",
			Name:     "Alice Smith",
		}

		resp, result := makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Response: %v", http.StatusCreated, resp.StatusCode, result)
		}

		if email, ok := result["email"].(string); !ok || email != "alice@example.com" {
			t.Errorf("Expected email alice@example.com, got %v", result["email"])
		}
	})

	var accessToken string
	var refreshToken string
	t.Run("login with valid credentials", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Email:    "alice@example.com",
			Password: "SecurePassword123",
		}

		resp, result := makeRequest(t, "POST", "/api/auth/login", loginReq, nil, server.URL)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %v", http.StatusOK, resp.StatusCode, result)
		}

		if token, ok := result["access_token"].(string); !ok || token == "" {
			t.Error("Expected access_token in response")
		} else {
			accessToken = token
		}

		if token, ok := result["refresh_token"].(string); !ok || token == "" {
			t.Error("Expected refresh_token in response")
		} else {
			refreshToken = token
		}
	})

	t.Run("access protected endpoint with valid token", func(t *testing.T) {
		resp, result := makeRequest(t, "GET", "/api/auth/profile", nil, &accessToken, server.URL)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %v", http.StatusOK, resp.StatusCode, result)
		}

		if user, ok := result["user"]; ok {
			if userMap, ok := user.(map[string]interface{}); ok {
				if email, ok := userMap["email"].(string); !ok || email != "alice@example.com" {
					t.Errorf("Expected email alice@example.com, got %v", userMap["email"])
				}
			}
		}
	})

	t.Run("refresh access token", func(t *testing.T) {
		refreshReq := dto.RefreshTokenRequest{
			RefreshToken: refreshToken,
		}

		resp, result := makeRequest(t, "POST", "/api/auth/refresh", refreshReq, nil, server.URL)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %v", http.StatusOK, resp.StatusCode, result)
		}

		if token, ok := result["access_token"].(string); !ok || token == "" {
			t.Error("Expected new access_token in response")
		} else {
			accessToken = token
		}
	})

	t.Run("update profile with valid token", func(t *testing.T) {
		updateReq := dto.UpdateProfileRequest{
			Name: "Alice Johnson",
		}

		resp, result := makeRequest(t, "PUT", "/api/auth/profile", updateReq, &accessToken, server.URL)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Response: %v", http.StatusOK, resp.StatusCode, result)
		}

		if user, ok := result["user"]; ok {
			if userMap, ok := user.(map[string]interface{}); ok {
				if name, ok := userMap["name"].(string); !ok || name != "Alice Johnson" {
					t.Errorf("Expected name Alice Johnson, got %v", userMap["name"])
				}
			}
		}
	})
}

func TestDuplicateRegistration(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	registerReq := dto.RegisterRequest{
		Email:    "bob@example.com",
		Password: "SecurePassword123",
		Name:     "Bob Smith",
	}

	resp, _ := makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	resp, result := makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status %d for duplicate, got %d. Response: %v", http.StatusConflict, resp.StatusCode, result)
	}
}

func TestInvalidCredentials(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	registerReq := dto.RegisterRequest{
		Email:    "charlie@example.com",
		Password: "SecurePassword123",
		Name:     "Charlie Smith",
	}
	makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)

	loginReq := dto.LoginRequest{
		Email:    "charlie@example.com",
		Password: "WrongPassword",
	}

	resp, result := makeRequest(t, "POST", "/api/auth/login", loginReq, nil, server.URL)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %d for invalid credentials, got %d. Response: %v", http.StatusUnauthorized, resp.StatusCode, result)
	}
}

func TestProtectedEndpointWithoutToken(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	resp, result := makeRequest(t, "GET", "/api/auth/profile", nil, nil, server.URL)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %d without token, got %d. Response: %v", http.StatusUnauthorized, resp.StatusCode, result)
	}
}

func TestInputValidation(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	tests := []struct {
		name           string
		registerReq    dto.RegisterRequest
		expectedStatus int
	}{
		{
			name: "invalid email format",
			registerReq: dto.RegisterRequest{
				Email:    "notanemail",
				Password: "SecurePassword123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			registerReq: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "short",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "name too short",
			registerReq: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "SecurePassword123",
				Name:     "A",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			registerReq: dto.RegisterRequest{
				Email:    "",
				Password: "SecurePassword123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, result := makeRequest(t, "POST", "/api/auth/register", tt.registerReq, nil, server.URL)
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %v", tt.expectedStatus, resp.StatusCode, result)
			}
		})
	}
}

func TestAdminFlow(t *testing.T) {
	server, f, cleanup := setupTestServer(t)
	defer cleanup()

	adminUser := "admin@example.com"
	adminPassword := "AdminPassword123"
	adminName := "Admin User"

	registerReq := dto.RegisterRequest{
		Email:    adminUser,
		Password: adminPassword,
		Name:     adminName,
	}

	resp, _ := makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create admin user: status %d", resp.StatusCode)
	}

	_, err := f.DB.Exec("UPDATE users SET role = 'admin' WHERE email = ?", adminUser)
	if err != nil {
		t.Fatalf("Failed to update user role: %v", err)
	}

	loginReq := dto.LoginRequest{
		Email:    adminUser,
		Password: adminPassword,
	}

	resp, result := makeRequest(t, "POST", "/api/auth/login", loginReq, nil, server.URL)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to login as admin: status %d, response %v", resp.StatusCode, result)
	}

	var adminToken string
	if token, ok := result["access_token"].(string); ok {
		adminToken = token
	} else {
		t.Fatalf("No access token in login response")
	}

	resp, result = makeRequest(t, "GET", "/api/admin/dashboard", nil, &adminToken, server.URL)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d for admin endpoint, got %d. Response: %v", http.StatusOK, resp.StatusCode, result)
	}
}

func TestNonAdminCannotAccessAdminEndpoint(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	registerReq := dto.RegisterRequest{
		Email:    "regular@example.com",
		Password: "RegularPassword123",
		Name:     "Regular User",
	}
	makeRequest(t, "POST", "/api/auth/register", registerReq, nil, server.URL)

	loginReq := dto.LoginRequest{
		Email:    "regular@example.com",
		Password: "RegularPassword123",
	}

	resp, result := makeRequest(t, "POST", "/api/auth/login", loginReq, nil, server.URL)
	var token string
	if t, ok := result["access_token"].(string); ok {
		token = t
	}

	resp, result = makeRequest(t, "GET", "/api/admin/dashboard", nil, &token, server.URL)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status %d for non-admin access, got %d. Response: %v", http.StatusForbidden, resp.StatusCode, result)
	}
}
