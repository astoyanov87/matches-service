package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Match struct {
	MatchId         string `json:"matchID"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Round           string `json:"round"`
	HomePlayerScore int    `json:"homePlayerScore"`
	AwayPlayerScore int    `json:"awayPlayerScore"`
}

type SortedByRound []Match

// Initialize Redis connection
func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.133.66.119:6379", // Redis address
		Password: "",                   // No password
		DB:       0,                    // Default DB
	})

	// Ping Redis to test connection
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return rdb
}

// Fetch matches by status (stored in Redis sets)
func getMatchesByStatus(ctx context.Context, rdb *redis.Client, status string) ([]Match, error) {
	// Assuming matches are stored in Redis sets by status (e.g., "matches:live", "matches:completed")

	var foundMatches []Match

	setKey := status + "_matches"
	fmt.Println("Requested: " + setKey)
	matchesIDs, err := rdb.SMembers(ctx, setKey).Result()
	if err != nil {
		return nil, err
	}

	for _, foundMatchId := range matchesIDs {
		matchData, err := rdb.HGet(context.Background(), "match:"+foundMatchId, "data").Result()
		if err != nil {
			fmt.Println("Error retrieving match from cache")
		}
		var matchObj Match
		err = json.Unmarshal([]byte(matchData), &matchObj)
		if err != nil {
			fmt.Println("Error unmarshaling data from Redis")
		}

		//@TODO:need to find the reason why sets by status in redis get messed up
		if strings.EqualFold(matchObj.Status, status) { //this is a workaround to filter out the matches with the wrong statuses
			foundMatches = append(foundMatches, matchObj)
		}
	}
	sort.Sort(SortedByRound(foundMatches))
	return foundMatches, err
}

// Get All matches from Redis
func getAllMatches(ctx context.Context, rdb *redis.Client) ([]Match, error) {

	var cursor uint64
	var keys []string
	var err error
	var foundMatches []Match

	matchPrefix := "match:*"
	keys, cursor, err = rdb.Scan(ctx, cursor, matchPrefix, 300).Result()

	if err != nil {
		fmt.Printf("Can not find keys that start with %v", matchPrefix)
		return foundMatches, nil
	}

	for _, matchKey := range keys {

		matchData, err := rdb.HGet(context.Background(), matchKey, "data").Result()
		if err != nil {
			fmt.Println("Error retrieving match from cache")
			return foundMatches, nil
		}
		var matchObj Match
		err = json.Unmarshal([]byte(matchData), &matchObj)

		if err != nil {
			fmt.Println("Error unmarshaling data from Redis")
		}
		//fmt.Println(matchObj)
		foundMatches = append(foundMatches, matchObj)
	}
	sort.Sort(SortedByRound(foundMatches))

	return foundMatches, nil
}

// Len is the number of elements in the collection.
func (a SortedByRound) Len() int {
	return len(a)
}

// Swap swaps the elements with indexes i and j.
func (a SortedByRound) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less reports whether the element with index i should sort before the element with index j.
func (a SortedByRound) Less(i, j int) bool {
	return a[i].Round < a[j].Round
}
