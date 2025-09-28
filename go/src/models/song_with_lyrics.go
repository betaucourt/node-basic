package models

// SongWithLyrics represents a song with its complete lyrics
type SongWithLyrics struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Lyrics []string `json:"lyrics"`
}
