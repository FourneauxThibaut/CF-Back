package handlers

import (
	"net/http"

	"github.com/FourneauxThibaut/CF-Back/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	sqlcgen "github.com/FourneauxThibaut/CF-Back/internal/db/sqlc"
)

// Me handles GET /api/me. Returns the current user from the auth context.
func Me(c *gin.Context) {
	u, ok := auth.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// Profile returns a handler for GET /api/profile. It uses the auth user ID to load the profile from the database.
func Profile(queries *sqlcgen.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		parsed, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		var id pgtype.UUID
		id.Bytes = parsed
		id.Valid = true
		profile, err := queries.GetProfileByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}
