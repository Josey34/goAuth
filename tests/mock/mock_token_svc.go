package mock_test

type MockTokenService struct {
	GeneratedAccessTokens  []string
	GeneratedRefreshTokens []string
	ValidateShouldFail     bool
	ValidateError          error
}

func NewMockTokenService() *MockTokenService {
	return &MockTokenService{}
}

func (m *MockTokenService) GenerateAccess(userID, role string) (string, error) {
	token := "access-token-" + userID
	m.GeneratedAccessTokens = append(m.GeneratedAccessTokens, token)
	return token, nil
}

func (m *MockTokenService) GenerateRefresh(userID, role string) (string, error) {
	token := "refresh-token-" + userID
	m.GeneratedRefreshTokens = append(m.GeneratedRefreshTokens, token)
	return token, nil
}

func (m *MockTokenService) Validate(tokenString string) (map[string]interface{}, error) {
	if m.ValidateShouldFail {
		return nil, m.ValidateError
	}

	return map[string]interface{}{
		"sub": "test-user-id",
	}, nil
}
