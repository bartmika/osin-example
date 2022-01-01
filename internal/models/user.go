package models

import (
	"context"
	"time"
)

const (
	UserActiveState         = 1
	UserInactiveState       = 0
	UserSystemAdminRoleID   = 1
	UserTenantAdminRoleID   = 2
	UserTenantStudentRoleID = 3
)

type User struct {
	ID                uint64    `json:"id,omitempty"`
	UUID              string    `json:"uuid,omitempty"`
	TenantID          uint64    `json:"tenant_id,omitempty"`
	Email             string    `json:"email,omitempty"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	Name              string    `json:"name,omitempty"`
	LexicalName       string    `json:"lexical_name,omitempty"`
	PasswordAlgorithm string    `json:"password_algorithm,omitempty"`
	PasswordHash      string    `json:"password_hash,omitempty"`
	State             int8      `json:"state,omitempty"`
	RoleID            int8      `json:"role_id,omitempty"`
	Timezone          string    `json:"timezone,omitempty"`
	Language          string    `json:"language"`
	CreatedTime       time.Time `json:"created_time,omitempty"`
	ModifiedTime      time.Time `json:"modified_time,omitempty"`
	JoinedTime        time.Time `json:"joined_time,omitempty"`
	Salt              string    `json:"salt,omitempty"`
	WasEmailActivated bool      `json:"was_email_activated,omitempty"`
	PrAccessCode      string    `json:"pr_access_code,omitempty"`
	PrExpiryTime      time.Time `json:"pr_expiry_time,omitempty"`
	AccessToken       string    `json:"access_token,omitempty"`
	RefreshToken      string    `json:"refresh_token,omitempty"`
}

type UserLite struct {
	ID       uint64 `json:"id,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	TenantID uint64 `json:"tenant_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	State    int8   `json:"state,omitempty"`
	RoleID   int8   `json:"role_id,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Language string `json:"language"`
}

type UserRepository interface {
	Insert(ctx context.Context, u *User) error
	UpdateByID(ctx context.Context, u *User) error
	UpdateByEmail(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUUID(ctx context.Context, uid string) (*User, error)
	CheckIfExistsByID(ctx context.Context, id uint64) (bool, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	InsertOrUpdateByID(ctx context.Context, u *User) error
	InsertOrUpdateByEmail(ctx context.Context, u *User) error
}
