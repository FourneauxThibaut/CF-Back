package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents the user info returned by Supabase Auth.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// SupabaseUserResponse is the response shape of GET /auth/v1/user.
type SupabaseUserResponse struct {
	ID           string            `json:"id"`
	Email        string            `json:"email"`
	UserMetadata map[string]any    `json:"user_metadata,omitempty"`
}

const (
	userContextKey  = "user"
	userIDContextKey = "userID"
)

// ValidateSupabaseToken returns a Gin middleware that verifies the JWT with Supabase Auth.
func ValidateSupabaseToken(supabaseURL, anonKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		token := parts[1]

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, strings.TrimSuffix(supabaseURL, "/")+"/auth/v1/user", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth request"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", anonKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Auth service unavailable"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "detail": string(body)})
			c.Abort()
			return
		}

		var userResp SupabaseUserResponse
		if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user"})
			c.Abort()
			return
		}

		c.Set(userContextKey, User{ID: userResp.ID, Email: userResp.Email})
		c.Set(userIDContextKey, userResp.ID)
		c.Next()
	}
}

// GetUser returns the user from the request context (set by ValidateSupabaseToken).
func GetUser(c *gin.Context) (User, bool) {
	u, ok := c.Get(userContextKey)
	if !ok {
		return User{}, false
	}
	user, ok := u.(User)
	return user, ok
}

// GetUserID returns the current user's ID or an empty string.
func GetUserID(c *gin.Context) string {
	id, _ := c.Get(userIDContextKey)
	s, _ := id.(string)
	return s
}
