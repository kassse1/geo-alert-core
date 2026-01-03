package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rqkhm/geo-alert-core/internal/config"
	"github.com/rqkhm/geo-alert-core/internal/handler"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()
	fmt.Println("Geo Alert Core is running on port", cfg.Port)

	// Регистрируем HTTP-эндпоинты
	http.HandleFunc("/health", handler.HealthHandler)

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
