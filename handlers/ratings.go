package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// PageRating stores voting data for a single page
type PageRating struct {
	TotalVotes  int `json:"total_votes"`
	SumOfScores int `json:"sum_of_scores"`
}

// RatingStore manages the embedded memory map and JSON persistence
type RatingStore struct {
	mu      sync.RWMutex
	Ratings map[string]*PageRating `json:"ratings"`
	dbPath  string
}

// Global store instance
var Store *RatingStore

// InitRatings initializes the embedded memory store and loads from disk
func InitRatings() {
	dbPath := filepath.Join("data", "ratings.json")

	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Printf("Warning: Could not create data directory: %v", err)
	}

	Store = &RatingStore{
		Ratings: make(map[string]*PageRating),
		dbPath:  dbPath,
	}

	// Try to load existing data
	data, err := os.ReadFile(dbPath)
	if err == nil {
		if err := json.Unmarshal(data, &Store.Ratings); err != nil {
			log.Printf("Warning: Could not parse ratings.json: %v", err)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Warning: Could not read ratings.json: %v", err)
	}
}

// save writes the current map to disk. Caller must hold the mutex.
func (s *RatingStore) save() {
	data, err := json.MarshalIndent(s.Ratings, "", "  ")
	if err != nil {
		log.Printf("Error marshaling ratings: %v", err)
		return
	}
	if err := os.WriteFile(s.dbPath, data, 0644); err != nil {
		log.Printf("Error saving ratings to disk: %v", err)
	}
}

// GetScoreAndVotes returns the calculated percentage score and total vote count for a given page slug
func (s *RatingStore) GetScoreAndVotes(slug string) (int, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rating, exists := s.Ratings[slug]
	if !exists || rating.TotalVotes == 0 {
		return 10, 0 // Baseline default score, 0 votes
	}

	// Calculate percentage: (Sum / (TotalVotes * 5)) * 100
	// Example: 5 votes, all 5 stars = 25. (25 / 25) * 100 = 100%
	// Example: 1 vote, 1 star = 1. (1 / 5) * 100 = 20%
	percentage := float64(rating.SumOfScores) / float64(rating.TotalVotes*5) * 100.0
	return int(percentage), rating.TotalVotes
}

// AddVote records a new vote (1-5), persists it, and returns the new percentage and total votes
func (s *RatingStore) AddVote(slug string, score int) (int, int) {
	// Clamp score between 1 and 5
	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Ratings[slug]; !exists {
		s.Ratings[slug] = &PageRating{}
	}

	s.Ratings[slug].TotalVotes++
	s.Ratings[slug].SumOfScores += score

	// Persist to JSON
	s.save()

	// We need to return the new percentage without locking again
	percentage := float64(s.Ratings[slug].SumOfScores) / float64(s.Ratings[slug].TotalVotes*5) * 100.0
	return int(percentage), s.Ratings[slug].TotalVotes
}

// RateRequest represents the expected JSON payload for voting
type RateRequest struct {
	Slug  string `json:"slug"`
	Score int    `json:"score"`
}

// RateHandler handles incoming POST requests to submit a vote
func RateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Slug == "" || req.Score < 1 || req.Score > 5 {
		http.Error(w, "Invalid slug or score (must be 1-5)", http.StatusBadRequest)
		return
	}

	// Add vote and get new percentage
	newPercentage, totalVotes := Store.AddVote(req.Slug, req.Score)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"percentage":  newPercentage,
		"total_votes": totalVotes,
	})
}
