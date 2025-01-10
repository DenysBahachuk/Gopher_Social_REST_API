package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UsersStore struct {
	redisDb *redis.Client
}

const UserExpTime = time.Minute

func (s *UsersStore) Get(ctx context.Context, userId int64) (*store.User, error) {
	key := fmt.Sprintf("user_%v", userId)

	data, err := s.redisDb.Get(ctx, key).Result()
	if err != nil {
		switch err {
		case redis.Nil:
			return nil, nil
		default:
			return nil, err
		}
	}

	var user store.User

	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
func (s *UsersStore) Set(ctx context.Context, user *store.User) error {
	key := fmt.Sprintf("user_%v", user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.redisDb.SetEX(ctx, key, data, UserExpTime).Err()
}

func (s *UsersStore) Delete(ctx context.Context, userID int64) error {
	cacheKey := fmt.Sprintf("user_%d", userID)
	return s.redisDb.Del(ctx, cacheKey).Err()
}
