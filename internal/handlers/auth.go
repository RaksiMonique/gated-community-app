package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gated-community-api/internal/domain"
	"gated-community-api/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLogin(db *pgxpool.Pool, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "invalid request body"}, nil)
			return
		}

		var userID, hashedPassword string
		query := `SELECT id, password_hash FROM users WHERE email = $1 AND deleted_at IS NULL`
		err := db.QueryRow(r.Context(), query, input.Email).Scan(&userID, &hashedPassword)
		if err != nil {
			domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "invalid credentials"}, nil)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
			domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "invalid credentials"}, nil)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   userID,
			"email": input.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"iat":   time.Now().Unix(),
		})

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			domain.WriteJSON(w, http.StatusInternalServerError, domain.Envelope{Success: false, Message: "failed to generate token"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{
			Success: true,
			Data: map[string]string{
				"token": tokenString,
			},
		}, nil)
	}
}

func HandleGetMe(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser := r.Context().Value(middleware.UserContextKey).(*middleware.AuthUser)

		var user struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}

		query := `SELECT id, email, first_name, last_name FROM users WHERE id = $1`
		err := db.QueryRow(r.Context(), query, authUser.ID).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			domain.WriteJSON(w, http.StatusNotFound, domain.Envelope{Success: false, Message: "user not found"}, nil)
			return
		}

		domain.WriteJSON(w, http.StatusOK, domain.Envelope{
			Success: true,
			Data:    user,
		}, nil)
	}
}
