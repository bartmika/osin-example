package models

import (
	"context"
	"time"
)

const (
	AuthorizedApplicationPermissionGrantedState = 1
	AuthorizedApplicationPermissionDeniedState  = 0
)

type AuthorizedApplication struct {
	ID            uint64    `json:"id"`
	UUID          string    `json:"uuid"`
	TenantID      uint64    `json:"tenant_id"`
	ApplicationID uint64    `json:"application_id"`
	UserID        uint64    `json:"user_id"`
	SessionUUID   string    `json:"session_uuid"`
	State         int8      `json:"state"`
	CreatedTime   time.Time `json:"created_time"`
	// CreatedFromIP  string    `json:"created_from_ip"`
	ModifiedTime time.Time `json:"modified_time"`
	// ModifiedFromIP string    `json:"modified_from_ip"`
}

type AuthorizedApplicationRepository interface {
	Insert(ctx context.Context, u *AuthorizedApplication) error
	UpdateByID(ctx context.Context, u *AuthorizedApplication) error
	GetByID(ctx context.Context, id uint64) (*AuthorizedApplication, error)
	GetByUUID(ctx context.Context, uid string) (*AuthorizedApplication, error)
	GetByUserIDAndApplicationID(ctx context.Context, uid uint64, aid uint64) (*AuthorizedApplication, error)
	CheckIfExistsByID(ctx context.Context, id uint64) (bool, error)
	CheckIfPermissionGrantedByUserIDAndByApplicationID(ctx context.Context, uid uint64, aid uint64) (bool, error)
	InsertOrUpdateByID(ctx context.Context, u *AuthorizedApplication) error
	DeleteByID(ctx context.Context, id uint64) error
}
