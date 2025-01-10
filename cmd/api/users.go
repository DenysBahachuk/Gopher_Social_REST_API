package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userContext string

const userKey userContext = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	user, err := app.getUserById(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundErrorResponse(w, r, err)
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := app.writeResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

type FollowerUserPayload struct {
	UserId int64 `json:"userId"`
}

// FollowUser godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User followed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	err = app.store.Followers.Follow(r.Context(), followedID, followerUser.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			app.conflictErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeResponse(w, http.StatusNoContent, nil)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// UnfollowUser gdoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"User ID"
//	@Success		204		{string}	string	"User unfollowed"
//	@Failure		400		{object}	error	"User payload missing"
//	@Failure		404		{object}	error	"User not found"
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r)

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeResponse(w, http.StatusNoContent, nil)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// ActivateUser gdoc
//
//	@Summary		Activates/Registers a user
//	@Description	Activates/Registers a user by an invitation token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// func (app *application) userContextModdleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		userIdAsString := chi.URLParam(r, "userId")
// 		userId, err := strconv.ParseInt(userIdAsString, 10, 64)
// 		if err != nil {
// 			app.badRequestErrorResponse(w, r, err)
// 			return
// 		}

// 		user, err := app.store.Users.GetById(r.Context(), userId)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, store.ErrNotFound):
// 				app.notFoundErrorResponse(w, r, err)
// 			default:
// 				app.internalServerErrorResponse(w, r, err)
// 			}
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), userKey, user)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func (app *application) getUserFromContext(r *http.Request) *store.User {
	return r.Context().Value(userKey).(*store.User)
}
