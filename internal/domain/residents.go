package domain

import (
	"context"
	"time"
)

type Resident struct {
	ID               string     `json:"id"`
	CommunityID      string     `json:"community_id"`
	UnitID           string     `json:"unit_id"`
	HouseholdID      *string    `json:"household_id"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	Phone            string     `json:"phone"`
	ResidentType     string     `json:"resident_type"`
	Status           string     `json:"status"`
	MoveInDate       *time.Time `json:"move_in_date"`
	MoveOutDate      *time.Time `json:"move_out_date"`
	IsPrimaryContact bool       `json:"is_primary_contact"`
}

type ResidentRepository interface {
	CreateWithHousehold(ctx context.Context, res *Resident) error
	ListByCommunity(ctx context.Context, communityID string) ([]Resident, error)
	GetByID(ctx context.Context, id string) (*Resident, error)
	Update(ctx context.Context, res *Resident) error
	SetPrimaryContact(ctx context.Context, residentID, householdID string) error
}

type ResidentService interface {
	Onboard(ctx context.Context, res *Resident) error
	List(ctx context.Context, communityID string) ([]Resident, error)
}
