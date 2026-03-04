package service

type TokenService interface {
	GenerateAccess(userID, role string) (string, error)
	GenerateRefresh(userID, role string) (string, error)
	Validate(token string) (map[string]interface{}, error)
}
