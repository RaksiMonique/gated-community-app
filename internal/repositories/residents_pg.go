package repositories

import (
	"context"
	"gated-community-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresResidentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresResidentRepository(db *pgxpool.Pool) domain.ResidentRepository {
	return &postgresResidentRepository{db: db}
}

func (r *postgresResidentRepository) CreateWithHousehold(ctx context.Context, res *domain.Resident) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var householdID string
	// Ensure household exists for this unit
	err = tx.QueryRow(ctx,
		"INSERT INTO households (community_id, unit_id, name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id",
		res.CommunityID, res.UnitID, res.LastName+" Household").Scan(&householdID)

	if err != nil {
		err = tx.QueryRow(ctx, "SELECT id FROM households WHERE unit_id = $1", res.UnitID).Scan(&householdID)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO residents (community_id, unit_id, household_id, first_name, last_name, email, phone, resident_type, status, move_in_date, is_primary_contact)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err = tx.QueryRow(ctx, query,
		res.CommunityID, res.UnitID, householdID, res.FirstName, res.LastName,
		res.Email, res.Phone, res.ResidentType, res.Status, res.MoveInDate, res.IsPrimaryContact,
	).Scan(&res.ID)

	if err != nil {
		return err
	}

	res.HouseholdID = &householdID
	return tx.Commit(ctx)
}

func (r *postgresResidentRepository) ListByCommunity(ctx context.Context, communityID string) ([]domain.Resident, error) {
	query := `SELECT id, unit_id, first_name, last_name, email, phone, resident_type, status, move_in_date FROM residents WHERE community_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Query(ctx, query, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var residents []domain.Resident
	for rows.Next() {
		var res domain.Resident
		err := rows.Scan(&res.ID, &res.UnitID, &res.FirstName, &res.LastName, &res.Email, &res.Phone, &res.ResidentType, &res.Status, &res.MoveInDate)
		if err != nil {
			continue
		}
		residents = append(residents, res)
	}
	return residents, nil
}

// Implement GetByID, Update, and SetPrimaryContact similarly...
func (r *postgresResidentRepository) GetByID(ctx context.Context, id string) (*domain.Resident, error) {
	return nil, nil
}
func (r *postgresResidentRepository) Update(ctx context.Context, res *domain.Resident) error {
	return nil
}
func (r *postgresResidentRepository) SetPrimaryContact(ctx context.Context, residentID, householdID string) error {
	return nil
}
