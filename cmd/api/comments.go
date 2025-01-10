package main

import (
	"net/http"

	"github.com/DenysBahachuk/gopher_social/internal/store"
)

type createCommentsPayload struct {
	Content string `json:"content" validate:"omitempty,max=100"`
}

func (app *application) createCommentsHandler(w http.ResponseWriter, r *http.Request) {
	commentsPayload := createCommentsPayload{}

	err := readJSON(w, r, &commentsPayload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(commentsPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	post := r.Context().Value(postKey).(*store.Post)

	comments := store.Comment{
		PostId:  post.ID,
		UserId:  post.UserID,
		Content: commentsPayload.Content,
	}

	err = app.store.Comments.Create(r.Context(), &comments)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	err = app.writeResponse(w, http.StatusCreated, comments)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
