package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUsersStore{},
	}
}

type MockUsersStore struct{}

func (m *MockUsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (m *MockUsersStore) GetById(ctx context.Context, id int64) (*User, error) {
	return &User{}, nil
}

func (m *MockUsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return &User{}, nil
}

func (m *MockUsersStore) CreateAndInvite(ctx context.Context, user *User, email string, exp time.Duration) error {
	return nil
}

func (m *MockUsersStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUsersStore) Delete(ctx context.Context, id int64) error {
	return nil
}
