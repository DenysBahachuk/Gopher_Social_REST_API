package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postContext string

const postKey postContext = "post"

type СreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		СreatePostPayload	true	"Post payload"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	postPayload := СreatePostPayload{}

	err := readJSON(w, r, &postPayload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	user := app.getUserFromContext(r)

	post := store.Post{
		Content: postPayload.Content,
		Title:   postPayload.Title,
		Tags:    postPayload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.writeResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// GetPost godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := app.getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	post.Comments = comments

	if err = app.writeResponse(w, http.StatusOK, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// DeletePost godoc
//
//	@Summary		Deletes a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object} string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "postId")

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Posts.DeleteById(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// UpdatePost godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	postPayload := UpdatePostPayload{}

	err := readJSON(w, r, &postPayload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	post := app.getPostFromCtx(r)

	if postPayload.Title != nil {
		post.Title = *postPayload.Title
	}

	if postPayload.Content != nil {
		post.Title = *postPayload.Content
	}

	err = app.store.Posts.UpdateById(r.Context(), post)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err = app.writeResponse(w, http.StatusOK, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "postId")

		postId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, postId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundErrorResponse(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postKey, post)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) getPostFromCtx(r *http.Request) *store.Post {
	return r.Context().Value(postKey).(*store.Post)
}
