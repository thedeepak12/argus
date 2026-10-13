package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/thedeepak12/argus/internal/auth"
	"github.com/thedeepak12/argus/internal/db"
)

type AuthHandler struct {
	queries *db.Queries
}

func getRefreshTokenExpiry() time.Duration {
	expiry := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if expiry == "" {
		return 7 * 24 * time.Hour
	}

	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return 7 * 24 * time.Hour
	}

	return duration
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{queries: queries}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user, err := h.queries.CreateUser(context.Background(), db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := h.queries.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	valid, err := auth.VerifyPassword(user.Password, req.Password)
	if err != nil || !valid {
		WriteError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.String())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}
	expiresAt := time.Now().Add(getRefreshTokenExpiry())

	_, err = h.queries.CreateRefreshToken(context.Background(), db.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	refreshToken, err := h.queries.GetRefreshTokenByToken(context.Background(), req.RefreshToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		WriteError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	accessToken, err := auth.GenerateAccessToken(refreshToken.UserID.String())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	h.queries.DeleteRefreshToken(context.Background(), req.RefreshToken)
	WriteJSON(w, http.StatusOK, nil)
}
