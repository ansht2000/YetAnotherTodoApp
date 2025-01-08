package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ansht2000/YetAnotherTodoApp/internal/auth"
	"github.com/ansht2000/YetAnotherTodoApp/internal/database"
	"github.com/google/uuid"
)

type returnValueUser struct {
	Id uuid.UUID `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func makeJWTAndRefreshToken(userID uuid.UUID, secretKey string) (string, string, error) {
	jwt, err := auth.MakeJWT(userID, secretKey, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", "", err
	}

	return jwt, refreshToken, nil
}

func (cfg *apiConfig) handlerLoginUser(resWriter http.ResponseWriter, req *http.Request) {
	type parametersLoginUser struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parametersLoginUser{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding login request data", err)
		return
	}

	user, err := cfg.db.GetUserFromEmail(req.Context(), params.Email)
	if err == sql.ErrNoRows {
		respondWithError(resWriter, http.StatusUnauthorized, "incorrect email or password", err)
		return
	} else if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error retrieving user data", err)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.PasswordHash); err != nil {
		respondWithError(resWriter, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	jwt, refreshToken, err := makeJWTAndRefreshToken(user.ID, cfg.secretKey)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error making authorization tokens", err)
		return
	}

	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
	}
	_, err = cfg.db.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	resVal := returnValueUser{
		Id: user.ID,
		Name: user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	resWriter.Header().Add("Set-Cookie", fmt.Sprintf("jwt=%s", jwt))
	resWriter.Header().Add("Set-Cookie", fmt.Sprintf("refresh=%s", refreshToken))
	respondWithJSON(resWriter, http.StatusOK, resVal)
}

func (cfg *apiConfig) handlerCreateUser(resWriter http.ResponseWriter, req *http.Request) {
	type parametersCreateUser struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	params := parametersCreateUser{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error decoding user request data", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating user password", err)
		return
	}

	createUserParams := database.CreateUserParams{
		Name: params.Name,
		Email: params.Email,
		PasswordHash: hashedPass,
	}

	user, err := cfg.db.CreateUser(req.Context(), createUserParams)
	if err != nil {
		respondWithError(resWriter, http.StatusInternalServerError, "error creating user", err)
	}

	resVal := returnValueUser{
		Id: user.ID,
		Name: user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(resWriter, http.StatusCreated, resVal)
}