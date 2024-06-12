package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTable(db *sql.DB) {
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

// StartDB checks if the database file exists and creates it if it does not.
func Connection(dbFile string) (*sql.DB, error) {
	// Check if the database file exists and create it if it does not
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		file, err := os.Create(dbFile)
		if err != nil {
			return nil, fmt.Errorf("could not create database file: %v", err)
		}
		file.Close()
		log.Printf("Database file created at %s\n", dbFile)
	} else if err != nil {
		return nil, fmt.Errorf("could not check database file: %v", err)
	} else {
		log.Printf("Database file already exists at %s\n", dbFile)
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	log.Println("Database connected")
	return db, nil
}
