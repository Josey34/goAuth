package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goauth/dto"
	"goauth/handler"
	mock_test "goauth/tests/mock"
	"goauth/usecase"

	"github.com/gin-gonic/gin"
)

func TestGetProfileHandler(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		setupUser      bool
		expectedStatus int
		expectUser     bool
	}{
		{
			name:           "get existing user profile",
			email:          "test@example.com",
			setupUser:      true,
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name:           "user not found",
			email:          "nonexistent@example.com",
			setupUser:      false,
			expectedStatus: http.StatusNotFound,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()

			if tt.setupUser {
				mockRepo.Insert(tt.email, "Test User", "hash")
			}

			userUsecase := usecase.NewUserUsecase(mockRepo)
			userHandler := handler.NewUserHandler(userUsecase)

			req := httptest.NewRequest("GET", "/profile", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.email) // Use email as ID

			userHandler.GetProfile(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectUser && w.Body.Len() == 0 {
				t.Errorf("expected response body but got empty")
			}
		})
	}
}

func TestUpdateProfileHandler(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		newName        string
		setupUser      bool
		expectedStatus int
		expectUser     bool
	}{
		{
			name:           "update existing user profile",
			email:          "test@example.com",
			newName:        "Updated Name",
			setupUser:      true,
			expectedStatus: http.StatusOK,
			expectUser:     true,
		},
		{
			name:           "user not found",
			email:          "nonexistent@example.com",
			newName:        "New Name",
			setupUser:      false,
			expectedStatus: http.StatusNotFound,
			expectUser:     false,
		},
		{
			name:           "empty name in request",
			email:          "test@example.com",
			newName:        "",
			setupUser:      true,
			expectedStatus: http.StatusBadRequest,
			expectUser:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()

			if tt.setupUser {
				mockRepo.Insert(tt.email, "Original Name", "hash")
			}

			userUsecase := usecase.NewUserUsecase(mockRepo)
			userHandler := handler.NewUserHandler(userUsecase)

			reqBody := dto.UpdateProfileRequest{Name: tt.newName}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.email) // Use email as ID

			userHandler.UpdateProfile(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectUser && w.Body.Len() == 0 {
				t.Errorf("expected response body but got empty")
			}
		})
	}
}

func TestAdminDashboardHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		expectedStatus int
	}{
		{
			name:           "admin dashboard access",
			userID:         "admin-user-id",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "admin dashboard with empty user id",
			userID:         "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()
			userUsecase := usecase.NewUserUsecase(mockRepo)
			userHandler := handler.NewUserHandler(userUsecase)

			req := httptest.NewRequest("GET", "/admin/dashboard", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set("userID", tt.userID)

			userHandler.AdminDashboard(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if w.Body.Len() == 0 {
				t.Errorf("expected response body but got empty")
			}
		})
	}
}
