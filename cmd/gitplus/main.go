package main

import (
	"log"
	"os"

	"github.com/iciwhite/gitplus/internal/auth"
	"github.com/iciwhite/gitplus/internal/config"
	"github.com/iciwhite/gitplus/internal/ui"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	authService := auth.NewOAuthService(cfg)
	if !authService.IsAuthenticated() {
		if err := authService.StartAuthFlow(); err != nil {
			log.Fatalf("Auth failed: %v", err)
		}
	}

	tui := ui.NewTUI(cfg, authService)
	if err := tui.Run(); err != nil {
		log.Fatalf("UI error: %v", err)
		os.Exit(1)
	}
}