package repositories

import (
	"context"
	"gated-community-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresAnnouncementRepository struct {
	db *pgxpool.Pool
}

func NewPostgresAnnouncementRepository(db *pgxpool.Pool) domain.AnnouncementRepository {
	return &postgresAnnouncementRepository{db: db}
}

func (r *postgresAnnouncementRepository) Create(ctx context.Context, a *domain.Announcement) error {
	query := `
		INSERT INTO announcements (community_id, author_id, title, content, audience_type, priority, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		a.CommunityID, a.AuthorID, a.Title, a.Content, a.AudienceType, a.Priority, a.ExpiresAt,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *postgresAnnouncementRepository) ListActive(ctx context.Context, communityID string) ([]domain.Announcement, error) {
	query := `
		SELECT id, author_id, title, content, audience_type, priority, expires_at, created_at
		FROM announcements
		WHERE community_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY priority DESC, created_at DESC`

	rows, err := r.db.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []domain.Announcement
	for rows.Next() {
		var a domain.Announcement
		err := rows.Scan(&a.ID, &a.AuthorID, &a.Title, &a.Content, &a.AudienceType, &a.Priority, &a.ExpiresAt, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, nil
}

func (r *postgresAnnouncementRepository) Delete(ctx context.Context, id, communityID string) error {
	query := `DELETE FROM announcements WHERE id = $1 AND community_id = $2`
	_, err := r.db.Exec(ctx, query, id, communityID)
	return err
}
