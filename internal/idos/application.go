package idos

import (
	"context"
	"time"

	"github.com/google/uuid"
	null "gopkg.in/guregu/null.v4"

	"github.com/bartmika/osin-example/internal/models"
)

type ApplicationLiteFilterIDO struct {
	TenantID  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder null.String `json:"sort_order"`
	SortField null.String `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"last_seen_id"`
	Limit     uint64      `json:"limit"`
}

type ApplicationLiteListResponseIDO struct {
	NextID  uint64                    `json:"next_id,omitempty"`
	Count   uint64                    `json:"count"`
	Results []*models.ApplicationLite `json:"results"`
}

func NewApplicationLiteListResponseIDO(arr []*models.ApplicationLite, count uint64) *ApplicationLiteListResponseIDO {
	// Calculate next id.
	var nextID uint64
	if len(arr) > 0 {
		lastRecord := arr[len(arr)-1]
		nextID = lastRecord.ID
	}

	res := &ApplicationLiteListResponseIDO{ // Return through HTTP.
		Count:   count,
		Results: arr,
		NextID:  nextID,
	}

	return res
}

type ApplicationCreateRequestIDO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
	Scope       string `json:"scope"`
	RedirectURL string `json:"redirect_url"`
	ImageURL    string `json:"image_url"`
}

func ApplicationCreateRequestUnmarshal(ctx context.Context, ido *ApplicationCreateRequestIDO) (m *models.Application, e error) {
	tenantID := uint64(ctx.Value("user_tenant_id").(uint64))

	//
	// Return unmarhsaled results.
	//

	return &models.Application{
		UUID:         uuid.NewString(),
		TenantID:     tenantID,
		Name:         ido.Name,
		Description:  ido.Description,
		WebsiteURL:   ido.WebsiteURL,
		Scope:        ido.Scope,
		RedirectURL:  ido.RedirectURL,
		ImageURL:     ido.ImageURL,
		CreatedTime:  time.Now().UTC(),
		ModifiedTime: time.Now().UTC(),
		State:        models.ApplicationRunningState,
	}, nil
}

type ApplicationResponseIDO struct {
	ID           uint64    `json:"id"`
	UUID         string    `json:"uuid"`
	TenantID     uint64    `json:"tenant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	WebsiteURL   string    `json:"website_url"`
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

func ApplicationResponseMarshal(m *models.Application) *ApplicationResponseIDO {
	return &ApplicationResponseIDO{
		ID:           m.ID,
		UUID:         m.UUID,
		TenantID:     m.TenantID,
		Name:         m.Name,
		Description:  m.Description,
		WebsiteURL:   m.WebsiteURL,
		Scope:        m.Scope,
		RedirectURL:  m.RedirectURL,
		ImageURL:     m.ImageURL,
		State:        m.State,
		ClientID:     m.ClientID,
		CreatedTime:  m.CreatedTime,
		ModifiedTime: m.ModifiedTime,
	}
}

type ApplicationUpdateRequestIDO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	WebsiteURL  string `json:"website_url"`
	Scope       string `json:"scope"`
	RedirectURL string `json:"redirect_url"`
	ImageURL    string `json:"image_url"`
}
