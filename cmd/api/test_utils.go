package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenysBahachuk/gopher_social/internal/auth"
	"github.com/DenysBahachuk/gopher_social/internal/ratelimiter"
	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/DenysBahachuk/gopher_social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	//logger := zap.Must(zap.NewProduction()).Sugar()
	//outputs no logs for tests
	logger := zap.NewNop().Sugar()

	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewCacheMockStore()
	testAuth := &auth.TestAuthenticator{}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		rateLimiter:   rateLimiter,
		config:        cfg,
	}
}

func executeRequest(mux http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	return rr
}

func checkresponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
