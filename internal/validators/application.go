package validators

import (
	"encoding/json"

	"github.com/bartmika/osin-example/internal/idos"
)

func ValidateApplicationCreateFromRequest(dirtyData *idos.ApplicationCreateRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.WebsiteURL == "" {
		e["website_url"] = "missing value"
	}
	if dirtyData.Scope == "" {
		e["scope"] = "missing value"
	}
	if dirtyData.RedirectURL == "" {
		e["redirect_url"] = "missing value"
	}
	if dirtyData.ImageURL == "" {
		e["image_url"] = "missing value"
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

func ValidateApplicationUpdateFromRequest(dirtyData *idos.ApplicationUpdateRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
	}
	if dirtyData.Description == "" {
		e["description"] = "missing value"
	}
	if dirtyData.WebsiteURL == "" {
		e["website_url"] = "missing value"
	}
	if dirtyData.Scope == "" {
		e["scope"] = "missing value"
	}
	if dirtyData.RedirectURL == "" {
		e["redirect_url"] = "missing value"
	}
	if dirtyData.ImageURL == "" {
		e["image_url"] = "missing value"
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
