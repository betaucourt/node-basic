package queries

import (
	"database/sql"
	"music-app/src/models"
)

// InsertLyric inserts a new lyric line into the database
func InsertLyric(db *sql.DB, songID, line int, text string) error {
	_, err := db.Exec("INSERT INTO lyrics (song, line, text) VALUES (?, ?, ?)", songID, line, text)
	return err
}

// GetLyricsBySongID retrieves all lyrics for a specific song
func GetLyricsBySongID(db *sql.DB, songID int) ([]models.Lyric, error) {
	query := "SELECT song, line, text FROM lyrics WHERE song = ? ORDER BY line"
	rows, err := db.Query(query, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lyrics []models.Lyric
	for rows.Next() {
		var lyric models.Lyric
		err := rows.Scan(&lyric.Song, &lyric.Line, &lyric.Text)
		if err != nil {
			return nil, err
		}
		lyrics = append(lyrics, lyric)
	}

	return lyrics, rows.Err()
}

// GetLyricsTextBySongID retrieves lyrics text as string slice for a specific song
func GetLyricsTextBySongID(db *sql.DB, songID int) ([]string, error) {
	query := "SELECT text FROM lyrics WHERE song = ? ORDER BY line"
	rows, err := db.Query(query, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lyricsText []string
	for rows.Next() {
		var text string
		err := rows.Scan(&text)
		if err != nil {
			return nil, err
		}
		lyricsText = append(lyricsText, text)
	}

	return lyricsText, rows.Err()
}
