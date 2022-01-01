package validators

import (
	"encoding/json"

	"github.com/bartmika/osin-example/internal/idos"
)

func ValidateTenantUpdateFromRequest(dirtyData *idos.TenantUpdateRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Name == "" {
		e["name"] = "missing value"
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

// func ValidateSubmissionUpdateFromRequest(dirtyData *idos.SubmissionUpdateRequestIDO) (bool, string) {
// 	e := make(map[string]string)
//
// 	if dirtyData.Name == "" {
// 		e["name"] = "missing value"
// 	}
// 	if dirtyData.Description == "" {
// 		e["description"] = "missing value"
// 	}
//
// 	if len(e) != 0 {
// 		b, err := json.Marshal(e)
// 		if err != nil { // Defensive code
// 			return false, err.Error()
// 		}
// 		return false, string(b)
// 	}
// 	return true, ""
// }
