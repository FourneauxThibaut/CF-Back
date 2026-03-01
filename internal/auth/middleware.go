package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// User représente les infos utilisateur retournées par Supabase Auth
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	// Ajouter d'autres champs selon besoin (phone, user_metadata, etc.)
}

// SupabaseUserResponse structure de la réponse GET /auth/v1/user
type SupabaseUserResponse struct {
	ID    string                 `json:"id"`
	Email string                 `json:"email"`
	UserMetadata map[string]any  `json:"user_metadata,omitempty"`
}

// ValidateSupabaseToken vérifie le JWT auprès de Supabase Auth (recommandé par Supabase)
func ValidateSupabaseToken(supabaseURL, anonKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Let CORS preflight through without auth (browser sends OPTIONS without Authorization)
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}
		token := parts[1]

		req, err := http.NewRequestWithContext(c.Context(), http.MethodGet, strings.TrimSuffix(supabaseURL, "/")+"/auth/v1/user", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create auth request",
			})
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", anonKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "Auth service unavailable",
			})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"detail": string(body),
			})
		}

		var userResp SupabaseUserResponse
		if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse user",
			})
		}

		// Stocker l'utilisateur dans le contexte Fiber pour les handlers
		c.Locals("user", User{
			ID:    userResp.ID,
			Email: userResp.Email,
		})
		c.Locals("userID", userResp.ID)
		return c.Next()
	}
}

// GetUser extrait l'utilisateur du contexte (à appeler après le middleware)
func GetUser(c *fiber.Ctx) (User, bool) {
	u, ok := c.Locals("user").(User)
	return u, ok
}

// GetUserID retourne l'ID utilisateur ou une chaîne vide
func GetUserID(c *fiber.Ctx) string {
	id, _ := c.Locals("userID").(string)
	return id
}

// RequireUser retourne une erreur 401 si pas d'utilisateur
func RequireUser(c *fiber.Ctx) (User, error) {
	u, ok := GetUser(c)
	if !ok {
		return User{}, fmt.Errorf("user not in context")
	}
	return u, nil
}
