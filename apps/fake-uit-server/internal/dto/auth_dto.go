package dto

import "github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/model"

type LoginRequest struct {
	StudentID  string `json:"student_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

type LoginResponse struct {
	StudentID    string `json:"student_id"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RememberedMe bool   `json:"remember_me"`
}

func NewLoginResponse(session model.AuthSession, rememberMe bool) LoginResponse {
	return LoginResponse{
		StudentID:    session.StudentID,
		AccessToken:  session.Token,
		TokenType:    session.TokenType,
		RememberedMe: rememberMe,
	}
}
