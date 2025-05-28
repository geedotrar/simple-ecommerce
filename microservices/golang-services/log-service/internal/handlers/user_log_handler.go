package handlers

import (
	"context"
	"encoding/json"
	"log-service/config"
	"log-service/internal/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateLog(w http.ResponseWriter, r *http.Request) {
	var logData models.Log
	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	go func() {
		_, err := config.LogCollection.InsertOne(context.Background(), logData)
		if err != nil {
			http.Error(w, "Failed to insert log", http.StatusInternalServerError)
			return
		}
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Log Saved",
	})
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.LogCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var logs []models.Log
	if err := cursor.All(ctx, &logs); err != nil {
		http.Error(w, "Failed to parse logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
