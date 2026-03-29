package services

import (
	"context"
	"gated-community-api/internal/domain"
)

type residentService struct {
	repo domain.ResidentRepository
}

func NewResidentService(repo domain.ResidentRepository) domain.ResidentService {
	return &residentService{repo: repo}
}

func (s *residentService) Onboard(ctx context.Context, res *domain.Resident) error {
	// Validation Logic
	vErrors := make(map[string]string)
	if res.FirstName == "" || res.LastName == "" {
		vErrors["name"] = "First and Last name are required"
	}
	if res.UnitID == "" {
		vErrors["unit_id"] = "Resident must be assigned to a unit"
	}
	if len(vErrors) > 0 {
		return domain.NewValidationError(vErrors)
	}

	// Defaults
	if res.Status == "" {
		res.Status = "ACTIVE"
	}
	if res.ResidentType == "" {
		res.ResidentType = "TENANT"
	}

	return s.repo.CreateWithHousehold(ctx, res)
}

func (s *residentService) List(ctx context.Context, communityID string) ([]domain.Resident, error) {
	return s.repo.ListByCommunity(ctx, communityID)
}
