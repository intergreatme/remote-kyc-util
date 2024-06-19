package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Connect to a SQLite database using the write ahead log (WAL)
func Connect(dbFile string) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	// Enable Write-Ahead Logging (WAL) mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to enable WAL mode: %v", err)
	}

	log.Println("Database connected")
	return db, nil
}

func CreateTables(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS TRANSACTIONS (
        id INTEGER PRIMARY KEY,
        origin_tx_id TEXT,
        order_number TEXT,
        company TEXT,
        config_id TEXT,
        tx_id TEXT,
        payload TEXT,
        response TEXT,
        errors TEXT,
        tx_status TEXT,
        time DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("could not create transactions table: %s\n", err)
	}
}
