package store

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func NewRolesStore(db *sql.DB) *RolesStore {
	return &RolesStore{db: db}
}

func (s *RolesStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `SELECT id, level, description 
		FROM roles 
		WHERE name = $1 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := Role{}

	err := s.db.QueryRowContext(ctx, query, roleName).Scan(
		&role.Id,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}
