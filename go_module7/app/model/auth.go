package model

import "time"

type RegisterRequest struct {
	NIM      string  `json:"nim" validate:"required,min=3,max=30,number"`
	Name     string  `json:"name" validate:"required,min=3,max=30"`
	Grade    float64 `json:"grade" validate:"required,min=0,max=4"`
	Password string  `json:"password" validate:"required,min=8,max=72,nospace,strongpassword"`
}

type LoginRequest struct {
	NIM      string `json:"nim"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginRequest) GetIdentifier() string {
	if l.NIM != "" {
		return l.NIM
	}
	return l.Username
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // in seconds (e.g. 900)
}

type RefreshToken struct {
	ID        int64      `json:"id"`
	UserID    int        `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type AuthUser struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
