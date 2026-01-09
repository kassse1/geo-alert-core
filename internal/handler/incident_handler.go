package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/kassse1/geo-alert-core/internal/domain"
	"github.com/kassse1/geo-alert-core/internal/service"
)

type IncidentHandler struct {
	service            *service.IncidentService
	statsWindowMinutes int
}

func NewIncidentHandler(
	service *service.IncidentService,
	statsWindowMinutes int,
) *IncidentHandler {
	return &IncidentHandler{
		service:            service,
		statsWindowMinutes: statsWindowMinutes,
	}
}

/*
=====================
CREATE
POST /api/v1/incidents
=====================
*/

type createIncidentRequest struct {
	Title   string  `json:"title"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM int     `json:"radius_m"`
}

func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createIncidentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.RadiusM <= 0 {
		http.Error(w, "invalid incident data", http.StatusBadRequest)
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
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(incident)
}

/*
=====================
LIST
GET /api/v1/incidents?page&limit
=====================
*/

func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	incidents, err := h.service.List(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if incidents == nil {
		incidents = []domain.Incident{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incidents)
}

/*
=====================
GET BY ID
GET /api/v1/incidents/{id}
=====================
*/

func (h *IncidentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/incidents/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	incident, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if incident == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(incident)
}

/*
=====================
UPDATE
PUT /api/v1/incidents/{id}
=====================
*/

func (h *IncidentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/incidents/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req createIncidentRequest
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
		Active:  true,
	}

	if err := h.service.Update(incident); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
=====================
DEACTIVATE
DELETE /api/v1/incidents/{id}
=====================
*/

func (h *IncidentHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/incidents/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Deactivate(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
