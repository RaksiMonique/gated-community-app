package domain

import (
	"context"
	"time"
)

type Announcement struct {
	ID           string     `json:"id"`
	CommunityID  string     `json:"community_id"`
	AuthorID     string     `json:"author_id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	AudienceType string     `json:"audience_type"`
	Priority     bool       `json:"priority"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type AnnouncementRepository interface {
	Create(ctx context.Context, a *Announcement) error
	ListActive(ctx context.Context, communityID string) ([]Announcement, error)
	Delete(ctx context.Context, id, communityID string) error
}

type AnnouncementService interface {
	Create(ctx context.Context, a *Announcement) error
	ListActive(ctx context.Context, communityID string) ([]Announcement, error)
	Delete(ctx context.Context, id, communityID string) error
}
