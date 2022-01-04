package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `ApplicationLite` model.
type ApplicationLiteFilter struct {
	TenantID  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type ApplicationLite struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Scope    string `json:"scope"`
	ImageURL string `json:"image_url"`
	State    int8   `json:"state"`
}

type ApplicationLiteRepository interface {
	ListByFilter(ctx context.Context, filter *ApplicationLiteFilter) ([]*ApplicationLite, error)
	CountByFilter(ctx context.Context, filter *ApplicationLiteFilter) (uint64, error)
}
