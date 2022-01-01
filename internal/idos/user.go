package idos

import (
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type UserIDO struct {
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

func NewUserIDO(m *models.User) *UserIDO {
	return &UserIDO{
		ID:                m.ID,
		UUID:              m.UUID,
		TenantID:          m.TenantID,
		Email:             m.Email,
		FirstName:         m.FirstName,
		LastName:          m.LastName,
		Name:              m.Name,
		LexicalName:       m.LexicalName,
		PasswordAlgorithm: m.PasswordAlgorithm,
		PasswordHash:      m.PasswordHash,
		State:             m.State,
		RoleID:            m.RoleID,
		Timezone:          m.Timezone,
		Language:          m.Language,
		CreatedTime:       m.CreatedTime,
		ModifiedTime:      m.ModifiedTime,
		JoinedTime:        m.JoinedTime,
		Salt:              m.Salt,
		WasEmailActivated: m.WasEmailActivated,
		PrAccessCode:      m.PrAccessCode,
		PrExpiryTime:      m.PrExpiryTime,
		AccessToken:       m.AccessToken,
		RefreshToken:      m.RefreshToken,
	}
}

type UserUpdateRequestIDO struct {
	Name string `json:"name"`
}
