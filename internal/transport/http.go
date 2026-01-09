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

	// ---------- Repositories ----------
	incidentRepo := repository.NewIncidentPostgresRepository(db.DB)
	checkRepo := repository.NewLocationCheckPostgresRepository(db.DB)

	// ---------- Services ----------
	incidentService := service.NewIncidentService(
		incidentRepo,
		checkRepo,
	)

	webhookService := service.NewWebhookService(cfg.WebhookURL)

	locationService := service.NewLocationService(
		incidentRepo,
		checkRepo,
		webhookService,
	)

	// ---------- Handlers ----------
	incidentHandler := handler.NewIncidentHandler(
		incidentService,
		cfg.StatsTimeWindowMinutes,
	)

	locationHandler := handler.NewLocationHandler(locationService)

	// ---------- Public ----------
	mux.HandleFunc("/api/v1/location/check", locationHandler.Check)
	mux.HandleFunc("/api/v1/system/health", handler.Health)

	// ---------- Incidents stats (MUST BE BEFORE /{id}) ----------
	mux.Handle(
		"/api/v1/incidents/stats",
		middleware.APIKeyMiddleware(
			cfg.APIKey,
			http.HandlerFunc(incidentHandler.Stats),
		),
	)

	// ---------- Incidents collection ----------
	mux.Handle(
		"/api/v1/incidents",
		middleware.APIKeyMiddleware(
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
		),
	)

	// ---------- Incidents by ID ----------
	mux.Handle(
		"/api/v1/incidents/",
		middleware.APIKeyMiddleware(
			cfg.APIKey,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					incidentHandler.GetByID(w, r)
				case http.MethodPut:
					incidentHandler.Update(w, r)
				case http.MethodDelete:
					incidentHandler.Deactivate(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	)

	return mux
}
