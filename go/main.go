package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/config"
	"github.com/intergreatme/remote-kyc-util/database"
	"github.com/intergreatme/remote-kyc-util/handlers"
	"github.com/intergreatme/remote-kyc-util/routes"
)

func main() {
	// Fetch flags
	dbFile := flag.String("dbfile", "transactions.sqlite3", "SQLite3 database file")
	configID := flag.String("config", "", "Customer config ID")
	password := flag.String("password", "", "Private key password")
	flag.Parse()

	c := handlers.Config{}

	if *configID == "" || *password == "" {
		configFile, err := config.ReadConfigFile("config.yaml")
		if err != nil {
			log.Fatalf("Config ID or password not provided and could not read from config file: %v", err)
		}

		c.ID = configFile.ConfigID
		c.PvtKeyPassword = configFile.Password
	} else {
		c.ID = *configID
		c.PvtKeyPassword = *password
	}

	if c.ID == "" || c.PvtKeyPassword == "" {
		log.Fatal("Config ID and password are required either as flags or in the config.yaml file.")
	}

	db, err := database.Connection(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create transactions table if it does not exists
	database.SetupTable(db)

	handler := handlers.NewHandler(db, c)

	// Set up routes and start server
	routes.Router(handler)

	// Check if IGM certificate already exist or download if not. Also verify that a private key of your own signature is present
	err = certs.FetchCertificates(handler.Config.ID)
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
