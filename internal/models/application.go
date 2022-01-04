package models

import (
	"context"
	"time"
)

const (
	ApplicationInactiveState      = 0
	ApplicationDraftState         = 1
	ApplicationPendingReviewState = 2
	ApplicationShutdownState      = 4
	ApplicationRunningState       = 5
)

type Application struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	TenantID     uint64    `json:"tenant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Scope        string    `json:"scope"`
	RedirectURL  string    `json:"redirect_url"`
	ImageURL     string    `json:"image_url"`
	State        int8      `json:"state"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	CreatedTime  time.Time `json:"created_time"`
	// CreatedFromIP  string    `json:"created_from_ip"`
	ModifiedTime time.Time `json:"modified_time"`
	// ModifiedFromIP string    `json:"modified_from_ip"`
}

type ApplicationRepository interface {
	Insert(ctx context.Context, u *Application) error
	UpdateByID(ctx context.Context, u *Application) error
	GetByID(ctx context.Context, id uint64) (*Application, error)
	GetByUUID(ctx context.Context, uid string) (*Application, error)
	GetByClientID(ctx context.Context, cid string) (*Application, error)
	CheckIfExistsByID(ctx context.Context, id uint64) (bool, error)
	CheckIfRunningByClientID(ctx context.Context, clientID string) (bool, error)
	InsertOrUpdateByID(ctx context.Context, u *Application) error
	DeleteByID(ctx context.Context, id uint64) error
}
