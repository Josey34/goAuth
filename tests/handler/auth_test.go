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
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid registration",
			requestBody: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "empty email",
			requestBody: dto.RegisterRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()
			mockToken := mock_test.NewMockTokenService()
			authUsecase := usecase.NewAuthUsecase(mockRepo, 10, mockToken)
			authHandler := handler.NewAuthHandler(authUsecase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Register(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
	}{
		{
			name:           "valid login",
			email:          "test@example.com",
			password:       "password123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			email:          "notfound@example.com",
			password:       "password123",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()
			mockToken := mock_test.NewMockTokenService()

			if tt.name == "valid login" {
				hash, _ := generateBcryptHash("password123", 10)
				mockRepo.Insert(tt.email, "Test User", hash)
			}

			authUsecase := usecase.NewAuthUsecase(mockRepo, 10, mockToken)
			authHandler := handler.NewAuthHandler(authUsecase)

			reqBody := dto.LoginRequest{Email: tt.email, Password: tt.password}
			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			authHandler.Login(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func generateBcryptHash(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}
