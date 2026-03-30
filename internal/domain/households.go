package domain

import (
	"context"
	"time"
)

type Household struct {
	ID          string    `json:"id"`
	CommunityID string    `json:"community_id"`
	UnitID      string    `json:"unit_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type HouseholdRepository interface {
	ListByCommunity(ctx context.Context, communityID string) ([]Household, error)
}

type HouseholdService interface {
	List(ctx context.Context, communityID string) ([]Household, error)
}
