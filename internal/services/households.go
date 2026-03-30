package services

import (
	"context"
	"gated-community-api/internal/domain"
)

type householdService struct {
	repo domain.HouseholdRepository
}

func NewHouseholdService(repo domain.HouseholdRepository) domain.HouseholdService {
	return &householdService{repo: repo}
}

func (s *householdService) List(ctx context.Context, communityID string) ([]domain.Household, error) {
	return s.repo.ListByCommunity(ctx, communityID)
}
