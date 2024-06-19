package certs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/intergreatme/certcrypto"
	"github.com/intergreatme/remote-kyc-util/config"
)

func CheckDirectory(dir string) error {
	// Check if the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
		return errors.New("could not create keys directory")
	}
	return nil
}

func FetchCertificates(cnf config.Configuration) error {
	// Check keys directory exists
	err := CheckDirectory(cnf.CertDir)
	if err != nil {
		return err
	}

	certFile := filepath.Join(cnf.CertDir, cnf.CompanyID+".pem")
	pfxFile := filepath.Join(cnf.CertDir, cnf.PFXFilename)

	// Check if the igm_certs.pfx file exists
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		// Download certificate if it does not exist
		uri := fmt.Sprintf("%s/signkey/%s", cnf.URL, cnf.CompanyID)
		err = certcrypto.DownloadCert(uri, certFile)
		if err != nil {
			return errors.New(err.Error())
		}
		log.Printf("Certificate downloaded and saved to %s\n", certFile)
	} else if err != nil {
		return fmt.Errorf("could not check certificate file: %v", err)
	}

	// Check if the certs.pfx file exists
	if _, err := os.Stat(pfxFile); os.IsNotExist(err) {
		log.Fatalf("Error: %s file not found. It needs to be added manually.\n", pfxFile)
	} else if err != nil {
		return fmt.Errorf("could not check %s file: %v", cnf.PFXFilename, err)
	}

	return nil
}
