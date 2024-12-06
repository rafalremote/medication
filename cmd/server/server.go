package main

import (
	"log"
	_ "medication/api/swagger"
	"net/http"
	"os"

	"medication/config"
	"medication/internal/db"
	"medication/internal/handlers"
	"medication/internal/services"
	"medication/pkg/logger"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log := logger.New()

	var dbInterface db.DBInterface
	switch cfg.DBDriver {
	case "postgres":
		dbInterface, err = db.NewPostgresDB(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize PostgresDB: %v", err)
		}
	default:
		log.Fatalf("Unsupported DB driver: %s", cfg.DBDriver)
	}
	defer func() {
		if err := dbInterface.Close(); err != nil {
			log.Errorf("Error closing database connection: %v", err)
		}
	}()

	medicationService := services.NewMedicationService(dbInterface)

	r := chi.NewRouter()

	env := os.Getenv("TARGET_RELEASE")

	if env == "DEV" {
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	r.Route("/medications", func(r chi.Router) {
		handlers.RegisterMedicationRoutes(r, medicationService, log)
	})

	log.Infof("Starting server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
