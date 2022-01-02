package validators

import (
	"encoding/json"

	"github.com/bartmika/osin-example/internal/idos"
)

func ValidateRegisterRequest(dirtyData *idos.RegisterRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if dirtyData.LastName == "" {
		e["last_name"] = "missing value"
	}
	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.Language == "" {
		e["language"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}

func ValidateLoginRequest(dirtyData *idos.LoginRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}

func ValidateRefreshTokenRequest(dirtyData *idos.RefreshTokenRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.GrantType == "" {
		e["grant_type"] = "missing value"
	}
	if dirtyData.RefreshToken == "" {
		e["refresh_token"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}
