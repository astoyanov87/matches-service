package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

// GetMatchesByStatusHandler retrieves matches by status
// @Summary Retrieve matches by status
// @Description Get matches filtered by their status (live, completed or scheduled) and sorted by round
// @Tags matches
// @Param status query string true "Match Status"
// @Success 200 {array} Match
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/matches [get]
func getMatchesByStatusHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	// Extract "status" query parameter from URL (e.g., ?status=live)
	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "Status query parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch matches by status from Redis
	ctx := context.Background()
	matches, err := getMatchesByStatus(ctx, rdb, status)
	if err != nil {
		log.Printf("Error retrieving matches from Redis: %v", err)
		http.Error(w, "Error retrieving matches", http.StatusInternalServerError)
		return
	}
	fmt.Println(matches)
	// Respond with matches as JSON
	w.Header().Set("Content-Type", "application/json")
	if len(matches) == 0 {
		w.Write([]byte("[]")) // Return empty array if no matches found
	} else {
		json.NewEncoder(w).Encode(matches)
	}
}
