package idos

import (
	"time"

	"github.com/bartmika/osin-example/internal/models"
)

type TenantIDO struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	Name         string    `json:"name"`
	State        int8      `json:"state"`
	Timezone     string    `json:"timestamp"`
	Language     string    `json:"language"`
	CreatedTime  time.Time `json:"created_time"`
	ModifiedTime time.Time `json:"modified_time"`
}

func NewTenantIDO(m *models.Tenant) *TenantIDO {
	return &TenantIDO{
		ID:           m.ID,
		UUID:         m.UUID,
		Name:         m.Name,
		State:        m.State,
		Timezone:     m.Timezone,
		Language:     m.Language,
		CreatedTime:  m.CreatedTime,
		ModifiedTime: m.ModifiedTime,
	}
}

type TenantUpdateRequestIDO struct {
	Name string `json:"name"`
}
