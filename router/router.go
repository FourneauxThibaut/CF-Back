package router

import (
	"strings"
	"time"

	"github.com/FourneauxThibaut/CF-Back/internal/auth"
	"github.com/FourneauxThibaut/CF-Back/internal/config"
	"github.com/FourneauxThibaut/CF-Back/handlers"
	sqlcgen "github.com/FourneauxThibaut/CF-Back/internal/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New builds the Gin engine with all routes and middleware.
func New(cfg *config.Config, queries *sqlcgen.Queries, ruleHandler *handlers.RuleSystemHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	allowedOrigins := getCORSOrigins(cfg.FrontendURL)
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, o := range allowedOrigins {
				if o == origin {
					return true
				}
			}
			return false
		},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	// Public routes
	r.GET("/health", handlers.Health)

	// Protected API group (Supabase Auth required)
	// Note: gin-contrib/cors handles OPTIONS preflight automatically
	api := r.Group("/api", auth.ValidateSupabaseToken(cfg.SupabaseURL, cfg.SupabaseAnonKey))
	api.GET("/me", handlers.Me)
	api.GET("/profile", handlers.Profile(queries))

	// Rule systems (explicit routes)
	api.GET("/rule-systems", ruleHandler.ListSystems)
	api.POST("/rule-systems", ruleHandler.CreateSystem)
	api.GET("/rule-systems/:systemId", ruleHandler.GetSystem)
	api.PUT("/rule-systems/:systemId", ruleHandler.UpdateSystem)
	api.DELETE("/rule-systems/:systemId", ruleHandler.DeleteSystem)

	// Rules (reorder before :ruleId to avoid "reorder" being captured)
	api.POST("/rule-systems/:systemId/rules", ruleHandler.CreateRule)
	api.PUT("/rule-systems/:systemId/rules/reorder", ruleHandler.ReorderRules)
	api.PUT("/rule-systems/:systemId/rules/:ruleId", ruleHandler.UpdateRule)
	api.DELETE("/rule-systems/:systemId/rules/:ruleId", ruleHandler.DeleteRule)

	// Blocks (reorder before :blockId)
	api.POST("/rule-systems/:systemId/rules/:ruleId/blocks", ruleHandler.AddBlock)
	api.PUT("/rule-systems/:systemId/rules/:ruleId/blocks/reorder", ruleHandler.ReorderBlocks)
	api.PUT("/rule-systems/:systemId/rules/:ruleId/blocks/:blockId", ruleHandler.UpdateBlock)
	api.DELETE("/rule-systems/:systemId/rules/:ruleId/blocks/:blockId", ruleHandler.DeleteBlock)

	// Block definitions
	api.GET("/rule-systems/:systemId/block-definitions", ruleHandler.GetBlockDefinitions)
	api.POST("/rule-systems/:systemId/block-definitions", ruleHandler.CreateBlockDefinition)
	api.PUT("/rule-systems/:systemId/block-definitions/:defId", ruleHandler.UpdateBlockDefinition)
	api.DELETE("/rule-systems/:systemId/block-definitions/:defId", ruleHandler.DeleteBlockDefinition)

	return r
}

func getCORSOrigins(env string) []string {
	if env == "" {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	parts := strings.Split(env, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return origins
}
