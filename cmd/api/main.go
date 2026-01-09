package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/kassse1/geo-alert-core/internal/config"
	"github.com/kassse1/geo-alert-core/internal/transport"
	"github.com/kassse1/geo-alert-core/pkg/postgres"
)

func main() {
	_ = godotenv.Load()

	// 1. Load config
	cfg := config.Load()

	// 2. Connect to PostgreSQL
	db, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal("postgres connection failed:", err)
	}
	defer db.Close()

	// 3. Create router
	router := transport.NewRouter(db, cfg)

	// 4. Start server
	log.Println("Server started on port", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, router))
}
