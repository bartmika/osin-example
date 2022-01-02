package idos

import "time"

// The struct used to represent the user's `register` POST request data.
type RegisterRequestIDO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Language  string `json:"language"`
}

// The struct used to represent the system's response when the `register` POST request was a success.
type RegisterResponseIDO struct {
	Message string `json:"message"`
}

// LoginRequest struct used to represent the user's `login` POST request data.
type LoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// The struct used to represent the system's response when the `login` POST request was a success.
type LoginResponseIDO struct {
	TenantID         uint64 `json:"tenant_id"`
	TenantSchemaName string `json:"tenant_schema_name"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	RoleID           int8   `json:"role_id"`
	Language         string `json:"language"`

	// https://pkg.go.dev/golang.org/x/oauth2#Token
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type,omitempty"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry,omitempty"`
}

// The struct used to represent the user's `refresh token` POST request data.
type RefreshTokenRequestIDO struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// The struct used to represent the system's response when the `refresh token` POST request was a success.
type RefreshTokenResponseIDO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GreetingRequestIDO struct {
	Name string `json:"name"`
}
type GreetingResponseIDO struct {
	Message string `json:"message"`
}
