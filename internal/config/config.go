package config

import (
	"os"
)

// Config contient la configuration de l'application
type Config struct {
	// Server
	Port string

	// Frontend (CORS)
	FrontendURL string

	// Supabase
	SupabaseURL     string
	SupabaseAnonKey string

	// Database (connection string Supabase)
	DatabaseURL string
}

// Load charge la config depuis l'environnement
func Load() (*Config, error) {
	_ = os.Setenv("TZ", "UTC")
	return &Config{
		Port:            getEnv("PORT", "3000"),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:5173"),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey: getEnv("SUPABASE_ANON_KEY", ""),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
