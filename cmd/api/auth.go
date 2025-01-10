package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/DenysBahachuk/gopher_social/internal/mailer"
	"github.com/DenysBahachuk/gopher_social/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	User  *store.User
	Token string
}

// RegisterUser godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	userPayload := RegisterUserPayload{}

	if err := readJSON(w, r, &userPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(userPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	user := store.User{
		Username: userPayload.Username,
		Email:    userPayload.Email,
		Role: store.Role{
			Name: "user",
		},
	}

	if err := user.Password.Set(userPayload.Password); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	ctx := r.Context()

	if err := app.store.Users.CreateAndInvite(ctx, &user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestErrorResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  &user,
		Token: plainToken,
	}

	isProdEnv := app.config.env == "production"

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	//send mail
	// log.Println("Sleeping for test before sending the mail")
	// time.Sleep(time.Second * 5)

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		//rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerErrorResponse(w, r, err)
		return
	}
	app.logger.Info("Email sent with status code %v", status)

	if err := app.writeResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// CreateToken godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	Token
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	//parse payload credentials
	userPayload := CreateUserTokenPayload{}

	if err := readJSON(w, r, &userPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	if err := Validate.Struct(userPayload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	//fetch the user (check if the user exists) from the payload
	user, err := app.store.Users.GetByEmail(r.Context(), userPayload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	//check the password
	if err := user.Password.Compare(userPayload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	//generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	//send it to the client
	if err := app.writeResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}
