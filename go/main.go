package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/database"
	"github.com/intergreatme/remote-kyc-util/routes"
)

func main() {
	// Fetch config and initialize flags
	dbFile := flag.String("dbfile", "transactions.sqlite3", "SQLite3 database file")
	configID := flag.String("config", "", "Customer config ID")
	flag.Parse()

	if *configID == "" {
		log.Fatal("Config ID is required. \n Usage: go run . config=<CONFIG_ID_HERE>")
	}

	db, err := database.Connection(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create transactions table if it does not exists
	database.SetupTable(db)

	// Set up routes and start server
	routes.Router()

	log.Println()
	// Load Certificates
	err = certs.LoadCertificates(*configID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

// Use IGM service to collect the public RSA key for appicable KYC configuration

// This will act as an allowlist creator
// It integrates with the IGM integration api to start a new transaction
// We will just build into the IGM Demo on Dev

// This does mean, that we will have to fetch secrets, like the the Config Key

// Things to cover:
// - Allowlist entry
// - Status API
// - Feedback API (Perhaps not its being deprecated?)
// - Completion API
// - Could also add the Golang code for validating an ID Number

//  Use selfsign service
// Sign allowlist
// Verify any response receive.
