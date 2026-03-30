package middleware

import (
	"context"
	"net/http"
	"strings"

	"gated-community-api/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserContextKey contextKey = "user"
const CommunityContextKey contextKey = "community_id"

type AuthUser struct {
	ID    string
	Email string
}

// Authenticate verifies the JWT token and adds the user to the request context
func Authenticate(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "missing authentication token"}, nil)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "invalid or expired token"}, nil)
				return
			}

			userID, okSub := claims["sub"].(string)
			email, okEmail := claims["email"].(string)

			if !okSub || !okEmail {
				domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "invalid token claims"}, nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, &AuthUser{ID: userID, Email: email})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TenantScope ensures the X-Community-ID is present and the user has access
func TenantScope(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*AuthUser)
			if !ok {
				domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "auth required"}, nil)
				return
			}

			communityID := r.Header.Get("X-Community-ID")
			if communityID == "" {
				domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "missing X-Community-ID"}, nil)
				return
			}

			// Verify membership in community_users
			var exists bool
			query := `SELECT EXISTS(SELECT 1 FROM community_users WHERE user_id = $1 AND community_id = $2)`
			db.QueryRow(r.Context(), query, user.ID, communityID).Scan(&exists)

			if !exists {
				domain.WriteJSON(w, http.StatusForbidden, domain.Envelope{Success: false, Message: "not a member of this community"}, nil)
				return
			}

			ctx := context.WithValue(r.Context(), CommunityContextKey, communityID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole checks if the authenticated user has a specific role in the requested community
func RequireRole(db *pgxpool.Pool, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*AuthUser)
			if !ok {
				domain.WriteJSON(w, http.StatusUnauthorized, domain.Envelope{Success: false, Message: "authentication required"}, nil)
				return
			}

			// Expecting community ID from context for multi-tenant scoping
			communityID, okComm := r.Context().Value(CommunityContextKey).(string)
			if !okComm || communityID == "" {
				domain.WriteJSON(w, http.StatusBadRequest, domain.Envelope{Success: false, Message: "missing community context"}, nil)
				return
			}

			var role string
			query := `SELECT role FROM community_users WHERE user_id = $1 AND community_id = $2`
			err := db.QueryRow(r.Context(), query, user.ID, communityID).Scan(&role)
			if err != nil {
				domain.WriteJSON(w, http.StatusForbidden, domain.Envelope{Success: false, Message: "access denied to this community"}, nil)
				return
			}

			roleAllowed := false
			for _, r := range roles {
				if strings.ToUpper(r) == strings.ToUpper(role) {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				domain.WriteJSON(w, http.StatusForbidden, domain.Envelope{Success: false, Message: "insufficient permissions"}, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
