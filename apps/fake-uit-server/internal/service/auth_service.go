package service

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type AuthService interface {
	Login(studentID string, password string) (model.AuthSession, bool)
	ValidateToken(token string) (model.AuthUser, bool)
}

type authService struct {
	credentials map[string]string
	tokens      map[string]string
}

func NewAuthService() AuthService {
	return &authService{
		credentials: map[string]string{
			"1": "password",
			"2": "password",
		},
		tokens: map[string]string{
			"fake-token-1": "1",
			"fake-token-2": "2",
		},
	}
}

func (s *authService) Login(studentID string, password string) (model.AuthSession, bool) {
	expectedPassword, ok := s.credentials[studentID]
	if !ok || expectedPassword != password {
		return model.AuthSession{}, false
	}

	return model.AuthSession{
		StudentID: studentID,
		Token:     "fake-token-" + studentID,
		TokenType: "Bearer",
	}, true
}

func (s *authService) ValidateToken(token string) (model.AuthUser, bool) {
	studentID, ok := s.tokens[token]
	if !ok {
		return model.AuthUser{}, false
	}

	return model.AuthUser{StudentID: studentID}, true
}
