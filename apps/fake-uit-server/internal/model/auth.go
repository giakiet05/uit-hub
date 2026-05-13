package model

type AuthSession struct {
	StudentID string
	Token     string
	TokenType string
}

type AuthUser struct {
	StudentID string
}
