package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FourneauxThibaut/CF-Back/internal/auth"
	"github.com/FourneauxThibaut/CF-Back/internal/config"
	"github.com/FourneauxThibaut/CF-Back/internal/db"
	sqlcgen "github.com/FourneauxThibaut/CF-Back/internal/db/sqlc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if cfg.DatabaseURL == "" || cfg.SupabaseURL == "" || cfg.SupabaseAnonKey == "" {
		log.Fatal("Missing required env: DATABASE_URL, SUPABASE_URL, SUPABASE_ANON_KEY")
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// Queries SQLC (généré par sqlc generate)
	queries := sqlcgen.New(pool)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})
	app.Use(recover.New())
	allowedOrigin := getCORSOrigins(cfg.FrontendURL)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	app.Use(logger.New())

	// Routes publiques
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// CORS preflight: handle OPTIONS /api/* with explicit headers so browser always gets Allow-Origin.
	app.Options("/api/*", func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin == "" {
			origin = allowedOrigin
		}
		if origin != allowedOrigin {
			return c.SendStatus(fiber.StatusForbidden)
		}
		c.Set("Access-Control-Allow-Origin", origin)
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "86400")
		return c.SendStatus(fiber.StatusNoContent)
	})

	// Groupe protégé par Supabase Auth
	api := app.Group("/api", auth.ValidateSupabaseToken(cfg.SupabaseURL, cfg.SupabaseAnonKey))
	api.Get("/me", func(c *fiber.Ctx) error {
		u, ok := auth.GetUser(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(u)
	})

	// Exemple: profil utilisateur (SQLC)
	api.Get("/profile", func(c *fiber.Ctx) error {
		userID := auth.GetUserID(c)
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		parsed, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}
		var id pgtype.UUID
		id.Bytes = parsed
		id.Valid = true
		profile, err := queries.GetProfileByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "profile not found"})
		}
		return c.JSON(profile)
	})

	addr := ":" + cfg.Port
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

// getCORSOrigins returns a slice of allowed origins (comma-separated FRONTEND_URL supported).
func getCORSOrigins(env string) string {
	if env == "" {
		return "http://localhost:5173"
	}
	return env
}
