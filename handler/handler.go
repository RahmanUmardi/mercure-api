package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"mercure-rocks/models"
)

type ChatHandler struct {
	HubURL         string
	PublisherToken string
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	msg.Timestamp = time.Now().Format(time.RFC3339)

	data, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to encode message", http.StatusInternalServerError)
		return
	}

	topic := "private:" + msg.Recipient
	if msg.GroupID != "" {
		topic = "group:" + msg.GroupID
	}

	payload := map[string]interface{}{
		"topic": topic,
		"data":  string(data),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode payload", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", h.HubURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+h.PublisherToken)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Authorization Header: %v", r.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to publish message, internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		http.Error(w, "Failed to publish message", resp.StatusCode)
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
