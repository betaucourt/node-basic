package database

import "database/sql"

// CreateTables creates the database tables for songs and lyrics
func CreateTables(db *sql.DB) error {
	// Create song table
	songTable := `
	CREATE TABLE IF NOT EXISTS song (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`

	// Create lyrics table
	lyricsTable := `
	CREATE TABLE IF NOT EXISTS lyrics (
		song INTEGER NOT NULL,
		line INTEGER NOT NULL,
		text TEXT NOT NULL,
		FOREIGN KEY (song) REFERENCES song(id),
		PRIMARY KEY (song, line)
	);`

	if _, err := db.Exec(songTable); err != nil {
		return err
	}

	if _, err := db.Exec(lyricsTable); err != nil {
		return err
	}

	return nil
}

// InsertSampleData inserts sample data into the database
func InsertSampleData(db *sql.DB) error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM song").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Data already exists
	}

	// Insert sample song
	result, err := db.Exec("INSERT INTO song (name) VALUES (?)", "Example Song")
	if err != nil {
		return err
	}

	songID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert sample lyrics
	lyrics := []string{
		"This is the first line of our example song",
		"Here comes the second line with a melody",
		"The third line continues the story",
		"And this is how our sample song ends",
	}

	for i, line := range lyrics {
		_, err := db.Exec("INSERT INTO lyrics (song, line, text) VALUES (?, ?, ?)", songID, i+1, line)
		if err != nil {
			return err
		}
	}

	return nil
}
