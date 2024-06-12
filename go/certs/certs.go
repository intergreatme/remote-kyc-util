package certs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/intergreatme/selfsign"
)

func CheckDirectory(dir string) error {
	// Check if the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
		return errors.New("could not create keys directory")
	}
	return nil
}

func FetchCertificates(configID string) error {
	keysDir := "keys"

	// Check keys directory exists
	err := CheckDirectory(keysDir)
	if err != nil {
		return err
	}

	certFile := filepath.Join(keysDir, "igm_certs.pem")
	keyFile := filepath.Join(keysDir, "key.pem")

	// Check if the igm_certs.pfx file exists
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		// Download certificate if it does not exist
		uri := fmt.Sprintf("https://dev.intergreatme.com/kyc/za/api/integration/signkey/%v", configID)
		err = selfsign.DownloadAndExtractCert(uri, keysDir, "igm_certs.pem")
		if err != nil {
			return fmt.Errorf("could not download key from IGM: %v", err)
		}
		log.Printf("Certificate downloaded and saved to %s\n", certFile)
	} else if err != nil {
		return fmt.Errorf("could not check cert file: %v", err)
	} else {
		log.Printf("Certificate already exists at %s\n", certFile)
	}

	// Check if the key.pem file exists
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Printf("Warning: key.pem file not found in %s. It needs to be added manually.\n", keysDir)
	} else if err != nil {
		return fmt.Errorf("could not check key file: %v", err)
	} else {
		log.Printf("Key file already exists at %s\n", keyFile)
	}

	return nil
}
