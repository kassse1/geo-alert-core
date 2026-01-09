package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)

		log.Println("WEBHOOK RECEIVED:")
		log.Printf("%+v\n", payload)

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Webhook mock listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
