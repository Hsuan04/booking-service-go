package booking

import (
	"database/sql"
)

type Repository interface {
	Save(studioID string, data []byte) (string, error)
	FindByID(id string) (*BookingDraft, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Save(studioID string, data []byte) (string, error) {
	var id string
	query := `INSERT INTO booking_drafts (studio_id, draft_data) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, studioID, data).Scan(&id)
	return id, err
}

func (r *postgresRepo) FindByID(id string) (*BookingDraft, error) {
	draft := &BookingDraft{}
	query := `SELECT id, studio_id, draft_data, created_at, expires_at 
	          FROM booking_drafts WHERE id = $1 AND expires_at > NOW()`

	err := r.db.QueryRow(query, id).Scan(
		&draft.ID, &draft.StudioID, &draft.DraftData, &draft.CreatedAt, &draft.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return draft, nil
}
