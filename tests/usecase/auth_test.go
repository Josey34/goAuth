package usecase_test

import (
	mock_test "goauth/tests/mock"
	"goauth/usecase"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		expectedError bool
		errorType     string
	}{
		{
			name:          "valid registration",
			email:         "new@example.com",
			password:      "secure123",
			expectedError: false,
		},
		{
			name:          "empty email",
			email:         "",
			password:      "secure123",
			expectedError: true,
			errorType:     "validation",
		},
		{
			name:          "duplicate email",
			email:         "existing@example.com",
			password:      "secure123",
			expectedError: true,
			errorType:     "conflict",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()

			if tt.name == "duplicate email" {
				mockRepo.Insert(tt.email, "Existing User", "hash")
			}

			authUsecase := usecase.NewAuthUsecase(mockRepo, 10, mock_test.NewMockTokenService())

			user, err := authUsecase.Register(tt.email, tt.password, "Test User")

			if tt.expectedError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedError && user == nil {
				t.Errorf("expected user but got nil")
			}
		})
	}
}
