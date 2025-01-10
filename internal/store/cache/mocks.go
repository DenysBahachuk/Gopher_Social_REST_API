package cache

import (
	"context"

	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/stretchr/testify/mock"
)

type MockUserCacheStore struct {
	mock.Mock
}

func NewCacheMockStore() Storage {
	return Storage{
		Users: &MockUserCacheStore{},
	}
}

func (m *MockUserCacheStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := m.Called(userID)
	return &store.User{}, args.Error(1)
}

func (m *MockUserCacheStore) Set(ctx context.Context, user *store.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserCacheStore) Delete(ctx context.Context, userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}
