package repositories

import (
	"context"
	"gated-community-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresHouseholdRepository struct {
	db *pgxpool.Pool
}

func NewPostgresHouseholdRepository(db *pgxpool.Pool) domain.HouseholdRepository {
	return &postgresHouseholdRepository{db: db}
}

func (r *postgresHouseholdRepository) ListByCommunity(ctx context.Context, communityID string) ([]domain.Household, error) {
	query := `
		SELECT id, community_id, unit_id, name, created_at 
		FROM households 
		WHERE community_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var households []domain.Household
	for rows.Next() {
		var h domain.Household
		err := rows.Scan(&h.ID, &h.CommunityID, &h.UnitID, &h.Name, &h.CreatedAt)
		if err != nil {
			return nil, err
		}
		households = append(households, h)
	}
	return households, nil
}
