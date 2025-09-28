package queries

import (
	"database/sql"
	"music-app/src/models"
)

// InsertSong inserts a new song into the database
func InsertSong(db *sql.DB, name string) (int64, error) {
	result, err := db.Exec("INSERT INTO song (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllSongs retrieves all songs from the database
func GetAllSongs(db *sql.DB) ([]models.Song, error) {
	query := "SELECT id, name FROM song ORDER BY id"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Name)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, rows.Err()
}

// GetSongByID retrieves a song by its ID
func GetSongByID(db *sql.DB, id int) (*models.Song, error) {
	query := "SELECT id, name FROM song WHERE id = ?"
	row := db.QueryRow(query, id)

	var song models.Song
	err := row.Scan(&song.ID, &song.Name)
	if err != nil {
		return nil, err
	}

	return &song, nil
}
