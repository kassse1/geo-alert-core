package transport

import (
	"net/http"

	"github.com/kassse1/geo-alert-core/internal/config"
	"github.com/kassse1/geo-alert-core/internal/handler"
	"github.com/kassse1/geo-alert-core/internal/middleware"
	"github.com/kassse1/geo-alert-core/internal/repository"
	"github.com/kassse1/geo-alert-core/internal/service"
	"github.com/kassse1/geo-alert-core/pkg/postgres"
)

func NewRouter(db *postgres.DB, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// =====================
	// Repositories
	// =====================
	incidentRepo := repository.NewIncidentPostgresRepository(db.DB)
	checkRepo := repository.NewLocationCheckPostgresRepository(db.DB)

	// =====================
	// Services
	// =====================
	incidentService := service.NewIncidentService(
		incidentRepo,
		checkRepo,
	)

	locationService := service.NewLocationService(
		incidentRepo,
		checkRepo,
	)

	// =====================
	// Handlers
	// =====================
	incidentHandler := handler.NewIncidentHandler(
		incidentService,
		cfg.StatsTimeWindowMinutes,
	)

	locationHandler := handler.NewLocationHandler(locationService)

	// =====================
	// Public endpoints
	// =====================
	mux.HandleFunc("/api/v1/system/health", handler.Health)
	mux.HandleFunc("/api/v1/location/check", locationHandler.Check)

	// =====================
	// Protected: incidents collection
	// =====================
	mux.Handle("/api/v1/incidents", middleware.APIKeyMiddleware(
		cfg.APIKey,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				incidentHandler.Create(w, r)
			case http.MethodGet:
				incidentHandler.List(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	// =====================
	// Protected: incidents by ID
	// =====================
	mux.Handle("/api/v1/incidents/", middleware.APIKeyMiddleware(
		cfg.APIKey,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				incidentHandler.Get(w, r)
			case http.MethodPut:
				incidentHandler.Update(w, r)
			case http.MethodDelete:
				incidentHandler.Delete(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	// =====================
	// Protected: stats
	// =====================
	mux.Handle("/api/v1/incidents/stats", middleware.APIKeyMiddleware(
		cfg.APIKey,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				incidentHandler.Stats(w, r)
				return
			}
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}),
	))

	return mux
}
