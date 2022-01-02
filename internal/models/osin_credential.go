package models

// UserData is used by the oAuth 2.0 system as extra data to include in the
// tokens that are sent back and forth by our system.
type UserData struct {
	TenantID uint64 `json:"tenant_id,omitempty"`
	UserID   uint64 `json:"user_id,omitempty"`
	UserUUID string `json:"user_uuid,omitempty"`
}
