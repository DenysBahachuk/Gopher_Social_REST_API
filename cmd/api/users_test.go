package main

import (
	"net/http"
	"testing"

	"github.com/DenysBahachuk/gopher_social/internal/store/cache"
	"github.com/stretchr/testify/mock"
)

func TestGetUserHandler(t *testing.T) {
	app := newTestApplication(t, config{})

	mux := app.mount()

	testToken, _ := app.authenticator.GenerateToken(nil)

	t.Run("should not allow unauthorized reqursts", func(t *testing.T) {
		//check for 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(mux, req)

		checkresponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(mux, req)

		checkresponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserCacheStore)

		app.config.redisCfg.enabled = true

		mockCacheStore.On("Get", int64(42)).Return(nil, nil)
		mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(mux, req)

		checkresponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.Calls = nil // Reset mock expectations
	})

	t.Run("should not hit the cache if it is not enabled", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserCacheStore)

		//app.config.redisCfg.enabled = false

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(mux, req)

		checkresponseCode(t, http.StatusOK, rr.Code)

		mockCacheStore.AssertNotCalled(t, "Get")

		mockCacheStore.Calls = nil // Reset mock expectations
	})
}
