package usecase_test

import (
	"testing"

	mock_test "goauth/tests/mock"
	"goauth/usecase"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupUser     bool
		expectedError bool
	}{
		{
			name:          "get existing user profile",
			email:         "test@example.com",
			setupUser:     true,
			expectedError: false,
		},
		{
			name:          "user not found",
			email:         "nonexistent@example.com",
			setupUser:     false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()

			if tt.setupUser {
				mockRepo.Insert(tt.email, "Test User", "hash")
			}

			userUsecase := usecase.NewUserUsecase(mockRepo)
			user, err := userUsecase.GetProfile(tt.email) // Use email as ID

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

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		newName       string
		setupUser     bool
		expectedError bool
		expectedName  string
	}{
		{
			name:          "update existing user",
			email:         "test@example.com",
			newName:       "Updated Name",
			setupUser:     true,
			expectedError: false,
			expectedName:  "Updated Name",
		},
		{
			name:          "update nonexistent user",
			email:         "nonexistent@example.com",
			newName:       "New Name",
			setupUser:     false,
			expectedError: true,
			expectedName:  "",
		},
		{
			name:          "update with empty name",
			email:         "test@example.com",
			newName:       "",
			setupUser:     true,
			expectedError: false,
			expectedName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_test.NewMockUserRepository()

			if tt.setupUser {
				mockRepo.Insert(tt.email, "Original Name", "hash")
			}

			userUsecase := usecase.NewUserUsecase(mockRepo)
			user, err := userUsecase.UpdateProfile(tt.email, tt.newName) // Use email as ID

			if tt.expectedError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedError && user == nil {
				t.Errorf("expected user but got nil")
			}
			if !tt.expectedError && user != nil && user.Name != tt.expectedName {
				t.Errorf("expected name %s but got %s", tt.expectedName, user.Name)
			}
		})
	}
}
