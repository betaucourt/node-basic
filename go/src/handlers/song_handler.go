package handlers

import (
	"database/sql"
	"encoding/json"
	"music-app/src/database/queries"
	"net/http"
)

// SongHandler handles song-related HTTP requests
type SongHandler struct {
	DB *sql.DB
}

// NewSongHandler creates a new SongHandler instance
func NewSongHandler(db *sql.DB) *SongHandler {
	return &SongHandler{DB: db}
}

// GetSongWithLyrics handles GET requests to retrieve a song with its lyrics
func (h *SongHandler) GetSongWithLyrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := queries.GetSongWithLyricsContext(r.Context(), h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
