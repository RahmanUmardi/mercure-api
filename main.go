package main

import (
	"log"
	"net/http"
	"time"

	"mercure-rocks/handler"

	"github.com/dunglas/mercure"
	"github.com/golang-jwt/jwt/v4"
)

func generateJWT(secret string, isPublisher bool) (string, error) {
	claims := jwt.MapClaims{
		"exp":               time.Now().Add(time.Hour * 24).Unix(),
		"mercure.publish":   []string{"*"},
		"mercure.subscribe": []string{"*"},
	}
	if !isPublisher {
		delete(claims, "mercure.publish")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func main() {
	secret := "rahman"

	publisherJWT, err := generateJWT(secret, true)
	if err != nil {
		log.Fatalf("Failed to generate publisher JWT: %v", err)
	}
	subscriberJWT, err := generateJWT(secret, false)
	if err != nil {
		log.Fatalf("Failed to generate subscriber JWT: %v", err)
	}

	log.Println("Publisher JWT:", publisherJWT)
	log.Println("Subscriber JWT:", subscriberJWT)

	hub, err := mercure.NewHub(
		mercure.WithPublisherJWT([]byte(secret), "HS256"),
		mercure.WithSubscriberJWT([]byte(secret), "HS256"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hub.Stop()

	chatHandler := &handler.ChatHandler{
		HubURL:         "http://localhost:8080/.well-known/mercure",
		PublisherToken: publisherJWT,
	}

	http.HandleFunc("/send-message", chatHandler.SendMessage)
	http.Handle("/.well-known/mercure", hub)

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
