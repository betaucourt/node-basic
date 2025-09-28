package queries

import (
	"context"
	"database/sql"
	"music-app/src/models"
)

// GetSongWithLyrics retrieves a song with its lyrics from the database
func GetSongWithLyrics(db *sql.DB) (*models.SongWithLyrics, error) {
	return GetSongWithLyricsContext(context.Background(), db)
}

// GetSongWithLyricsContext retrieves a song with its lyrics from the database with context
func GetSongWithLyricsContext(ctx context.Context, db *sql.DB) (*models.SongWithLyrics, error) {
	// Query to get song with lyrics
	query := `
	SELECT s.id, s.name, l.text 
	FROM song s 
	JOIN lyrics l ON s.id = l.song 
	ORDER BY s.id, l.line`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result models.SongWithLyrics
	var lyrics []string
	var currentSongID int

	for rows.Next() {
		var songID int
		var songName, lyricText string

		err := rows.Scan(&songID, &songName, &lyricText)
		if err != nil {
			return nil, err
		}

		if currentSongID == 0 {
			currentSongID = songID
			result.ID = songID
			result.Name = songName
		}

		lyrics = append(lyrics, lyricText)
	}

	result.Lyrics = lyrics

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
