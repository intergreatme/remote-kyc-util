package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/intergreatme/remote-kyc-util/certs"
	"github.com/intergreatme/remote-kyc-util/config"
	"github.com/intergreatme/remote-kyc-util/database"
	"github.com/intergreatme/remote-kyc-util/handlers"
	"github.com/intergreatme/remote-kyc-util/routes"

	_ "modernc.org/sqlite"
)

// getConfig provides an input-driven approach to configuring the application at start-up
func getConfig() config.Configuration {
	var cnf config.Configuration

	cnf, err := config.Read()
	// if we can read the config file, we can assume everything is fine and do not need to prompt
	// the user to add in the required information.
	if err == nil {
		return cnf
	}
	var input string
	fmt.Println("Config file not found, running interactively")

	fmt.Print("Company config ID: ")
	fmt.Scan(&input)
	cnf.CompanyID = input

	fmt.Print("URL to connect to: ")
	fmt.Scan(&input)
	cnf.URL = input

	fmt.Print("PFX filename: ")
	fmt.Scan(&input)
	cnf.PFXFilename = input

	fmt.Print("x509 Password: ")
	fmt.Scan(&input)
	cnf.Password = input

	cnf.CertDir = ".certs"
	fmt.Println("Please add your PFX files to the hidden directory .certs")

	err = cnf.Write()
	if err != nil {
		log.Fatal(err)
	}

	return cnf
}

const dbFile = "transactions.sqlite3"

func main() {
	cnf := getConfig()

	db, err := database.Connect(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create transactions table if it does not exists
	database.CreateTables(db)

	// Check if IGM certificate already exist or download if not.
	// Also verify that a private key of your own signature is present
	// If it does not exist, fail.
	err = certs.FetchCertificates(cnf)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler(db, cnf)

	// Set up routes and start server
	router := routes.Router(handler)

	log.Println("Starting server on http://localhost:8200/")
	if err := http.ListenAndServe(":8200", router); err != nil {
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
