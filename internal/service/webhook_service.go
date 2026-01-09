package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kassse1/geo-alert-core/internal/domain"
)

type WebhookService struct {
	url string
}

func NewWebhookService(url string) *WebhookService {
	return &WebhookService{url: url}
}

func (w *WebhookService) Send(userID string, incidents []domain.Incident) {
	if w.url == "" {
		return
	}

	payload := map[string]interface{}{
		"user_id":   userID,
		"incidents": incidents,
		"sent_at":   time.Now(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("webhook marshal error:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, w.url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("webhook request error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("webhook send error:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("webhook sent, status:", resp.StatusCode)
}
