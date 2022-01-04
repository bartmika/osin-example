package models

import (
	"context"
	// "time"
	null "gopkg.in/guregu/null.v4"
)

// Structure used to encapsulate the various filters we want to apply when we
// perform our `listing` functionality for the `AuthorizedApplicationLite` model.
type AuthorizedApplicationLiteFilter struct {
	TenantID  uint64      `json:"tenant_id"`
	States    []int8      `json:"states"`
	SortOrder string      `json:"sort_order"`
	SortField string      `json:"sort_field"`
	Search    null.String `json:"search"`
	Offset    uint64      `json:"offset"`
	Limit     uint64      `json:"limit"`
}

type AuthorizedApplicationLite struct {
	ID            uint64 `json:"id"`
	ApplicationID uint64 `json:"application_id"`
	State         int8   `json:"state"`
}

type AuthorizedApplicationLiteRepository interface {
	ListByFilter(ctx context.Context, filter *AuthorizedApplicationLiteFilter) ([]*AuthorizedApplicationLite, error)
	CountByFilter(ctx context.Context, filter *AuthorizedApplicationLiteFilter) (uint64, error)
}
