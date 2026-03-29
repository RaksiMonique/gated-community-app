package services

import (
	"context"
	"gated-community-api/internal/domain"
)

type announcementService struct {
	repo domain.AnnouncementRepository
}

func NewAnnouncementService(repo domain.AnnouncementRepository) domain.AnnouncementService {
	return &announcementService{repo: repo}
}

func (s *announcementService) Create(ctx context.Context, a *domain.Announcement) error {
	vErrors := make(map[string]string)
	if a.Title == "" {
		vErrors["title"] = "Title is required"
	}
	if len(a.Content) < 10 {
		vErrors["content"] = "Content must be at least 10 characters"
	}

	if len(vErrors) > 0 {
		return domain.NewValidationError(vErrors)
	}

	if a.AudienceType == "" {
		a.AudienceType = "ALL"
	}

	return s.repo.Create(ctx, a)
}

func (s *announcementService) ListActive(ctx context.Context, communityID string) ([]domain.Announcement, error) {
	return s.repo.ListActive(ctx, communityID)
}

func (s *announcementService) Delete(ctx context.Context, id, communityID string) error {
	return s.repo.Delete(ctx, id, communityID)
}
