package models

import (
	"context"
	"time"
)

const (
	TenantActiveState   = 1
	TenantInactiveState = 0
)

type Tenant struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	State        int8      `json:"state"`
	Timezone     string    `json:"timestamp"`
	Language     string    `json:"language"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

type TenantRepository interface {
	Insert(ctx context.Context, u *Tenant) error
	UpdateByID(ctx context.Context, u *Tenant) error
	GetByID(ctx context.Context, id uint64) (*Tenant, error)
	GetByName(ctx context.Context, name string) (*Tenant, error)
	CheckIfExistsByID(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByName(ctx context.Context, name string) (bool, error)
	InsertOrUpdateByID(ctx context.Context, u *Tenant) error
}
