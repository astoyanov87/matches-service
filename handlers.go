package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

// GetMatchesByStatusHandler retrieves matches by status
// @Summary Retrieve matches by status
// @Description Get matches filtered by their status (live, completed or scheduled) and sorted by round
// @Tags matches
// @Param status query string false "Match Status"
// @Success 200 {array} Match
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/matches [get]
func getMatchesByStatusHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {
	// Extract "status" query parameter from URL (e.g., ?status=live)
	status := r.URL.Query().Get("status")
	ctx := context.Background()
	fmt.Println("The status is:" + status)

	var returnMatches []Match
	if status == "" {
		fmt.Println("Requested all matches")
		matches, err := getAllMatches(ctx, rdb)

		if err != nil {
			log.Printf("Error retrieving matches from Redis: %v", err)
			http.Error(w, "Error retrieving matches", http.StatusInternalServerError)
			return
		}
		returnMatches = matches
	} else {
		// Fetch matches by status from Redis

		matches, err := getMatchesByStatus(ctx, rdb, status)
		if err != nil {
			log.Printf("Error retrieving matches from Redis: %v", err)
			http.Error(w, "Error retrieving matches", http.StatusInternalServerError)
			return
		}
		//fmt.Println(matches)
		returnMatches = matches
	}

	// Respond with matches as JSON
	w.Header().Set("Content-Type", "application/json")
	if len(returnMatches) == 0 {
		w.Write([]byte("[]")) // Return empty array if no matches found
	} else {
		json.NewEncoder(w).Encode(returnMatches)
	}
}

// GetMatchByIdHandler retrieves matches by status
// @Summary Retrieve match details by match ID
// @Description Get match details by given match ID
// @Tags match
// @Param id path string true "Match ID"
// @Success 200 {object} Match
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/match/{id} [get]
func getMatchByIdHandler(rdb *redis.Client, w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	fmt.Println("The URL is:" + r.URL.Path)
	matchId := strings.TrimPrefix(r.URL.Path, "/api/v1/match/")

	fmt.Println("Match ID:" + matchId)
	if matchId == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}
	foundMatch, err := getMatchById(ctx, rdb, matchId)
	if err != nil {
		fmt.Println("Cant not fetch match")
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(foundMatch)

}
