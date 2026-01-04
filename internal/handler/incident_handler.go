package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/kassse1/geo-alert-core/internal/domain"
	"github.com/kassse1/geo-alert-core/internal/service"
)

// =====================
// Handler
// =====================

type IncidentHandler struct {
	service             *service.IncidentService
	defaultStatsMinutes int
}

type statsResponse struct {
	UserCount int `json:"user_count"`
}

func NewIncidentHandler(
	service *service.IncidentService,
	defaultStatsMinutes int,
) *IncidentHandler {
	return &IncidentHandler{
		service:             service,
		defaultStatsMinutes: defaultStatsMinutes,
	}
}

func (h *IncidentHandler) Stats(w http.ResponseWriter, r *http.Request) {
	minutesStr := r.URL.Query().Get("minutes")

	minutes := h.defaultStatsMinutes
	if minutesStr != "" {
		if m, err := strconv.Atoi(minutesStr); err == nil && m > 0 {
			minutes = m
		}
	}

	count, err := h.service.Stats(minutes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(statsResponse{
		UserCount: count,
	})
}

// =====================
// Requests
// =====================

type createIncidentRequest struct {
	Title   string  `json:"title"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m"`
}

type updateIncidentRequest struct {
	Title   string  `json:"title"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m"`
	Active  bool    `json:"active"`
}

// =====================
// Handlers
// =====================

// POST /api/v1/incidents
func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createIncidentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RadiusM <= 0 {
		http.Error(w, "radius must be positive", http.StatusBadRequest)
		return
	}

	incident := &domain.Incident{
		Title:   req.Title,
		Lat:     req.Lat,
		Lon:     req.Lon,
		RadiusM: req.RadiusM,
		Active:  true,
	}

	if err := h.service.Create(incident); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

// GET /api/v1/incidents?page=1&limit=10
func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	incidents, err := h.service.List(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incidents)
}

// GET /api/v1/incidents/{id}
func (h *IncidentHandler) Get(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	incident, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

// PUT /api/v1/incidents/{id}
func (h *IncidentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	incident := &domain.Incident{
		ID:      id,
		Title:   req.Title,
		Lat:     req.Lat,
		Lon:     req.Lon,
		RadiusM: req.RadiusM,
		Active:  req.Active,
	}

	if err := h.service.Update(incident); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/incidents/{id} (deactivate)
func (h *IncidentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Deactivate(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
