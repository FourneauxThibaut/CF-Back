package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FourneauxThibaut/CF-Back/handlers"
	"github.com/FourneauxThibaut/CF-Back/internal/config"
	"github.com/FourneauxThibaut/CF-Back/internal/db"
	"github.com/FourneauxThibaut/CF-Back/internal/ruleeditor"
	"github.com/FourneauxThibaut/CF-Back/router"
	sqlcgen "github.com/FourneauxThibaut/CF-Back/internal/db/sqlc"
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
		log.Fatal(`Missing required env: DATABASE_URL, SUPABASE_URL, SUPABASE_ANON_KEY.
Set them in Scalingo: Dashboard > Your App > Environment, or:
  scalingo --app <app-name> env-set DATABASE_URL="..." SUPABASE_URL="..." SUPABASE_ANON_KEY="..."`)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	queries := sqlcgen.New(pool)
	ruleRepo := ruleeditor.NewJSONBRepository(pool)
	ruleSvc := ruleeditor.NewService(ruleRepo)
	ruleHandler := handlers.NewRuleSystemHandler(ruleSvc)
	engine := router.New(cfg, queries, ruleHandler)

	addr := ":" + cfg.Port
	srv := &http.Server{Addr: addr, Handler: engine}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
