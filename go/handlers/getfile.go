package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intergreatme/certcrypto"
	"github.com/intergreatme/remote-kyc-util/file"
	"github.com/intergreatme/remote-kyc-util/request"
	"github.com/intergreatme/remote-kyc-util/response"
)

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var fp file.FilePayload

	err = json.Unmarshal(body, &fp)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Unix()

	// Marshal the Payload struct to JSON
	jsonData, err := json.Marshal(fp)
	if err != nil {
		http.Error(w, "Error converting payload to JSON", http.StatusInternalServerError)
		return
	}

	pfxFile := filepath.Join(h.Config.CertDir, h.Config.PFXFilename)
	privateKey, _, err := certcrypto.ReadPKCS12(pfxFile, h.Config.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasign := string(jsonData) + fmt.Sprintf("%d", timestamp)

	signature, err := certcrypto.SignData(privateKey, []byte(datasign))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sig := base64.StdEncoding.EncodeToString(signature)
	reqBody := request.RequestPayload{
		Payload:   string(jsonData),
		Timestamp: timestamp,
		Signature: sig,
	}

	// Make the POST request to the Allowlist API
	resp, err := request.GetFileAPI(reqBody, h.Config)
	if err != nil {
		http.Error(w, "API call failed, "+err.Error(), http.StatusInternalServerError)
		return
	}

	payloadBody, err := io.ReadAll(&resp.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("\n\nAPI response: %s", payloadBody)

	if resp.StatusCode != 200 {
		var rw response.ResponseWrapper
		err = rw.FromJSON(payloadBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("\n\nPayload error response: %v", rw)
		// Write the API response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write(payloadBody)
		return
		// If you want to store the error response somewhere, do that here using the wrapper response.
		// e.g., {
		// 	"code": 500,
		// 	"message": "There was an error processing your request. It has been logged (ID f4f981f01de5a056)."
		// }
		// Where, rw.code = 500 and rw.message = "the message seen above"
	}

	// Response 200 OK, will be binary data, as in the blob file.
	// Check response content-type to assign correct file extension.
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type not found in the response", http.StatusBadRequest)
		return
	}

	mimeType := strings.Split(contentType, ";")[0]
	fileExtension := ".png"
	switch mimeType {
	case "image/jpeg":
		fileExtension = ".jpeg"
	case "application/pdf":
		fileExtension = ".pdf"
	}

	fmt.Printf("The file extension for MIME type '%s' is '%s'\n", mimeType, fileExtension)

	// Set file name, i.e. txID_file.png
	fileName := fmt.Sprintf("%s_%s_%s%s", fp.TxID, fp.DocumentType, fp.FileType, fileExtension)
	err = os.WriteFile(fileName, payloadBody, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("File saved as '%s'\n", fileName)

	// Write file to webapp
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(payloadBody)
}

func (h *Handler) GetLivelinessFile(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var f file.LivelinessPayload
	err = json.Unmarshal(body, &f)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Unix()

	// Marshal the Payload struct to JSON
	jsonData, err := json.Marshal(f)
	if err != nil {
		http.Error(w, "Error converting payload to JSON", http.StatusInternalServerError)
		return
	}

	pfxFile := filepath.Join(h.Config.CertDir, h.Config.PFXFilename)
	privateKey, _, err := certcrypto.ReadPKCS12(pfxFile, h.Config.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasign := string(jsonData) + fmt.Sprintf("%d", timestamp)

	signature, err := certcrypto.SignData(privateKey, []byte(datasign))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sig := base64.StdEncoding.EncodeToString(signature)
	reqBody := request.RequestPayload{
		Payload:   string(jsonData),
		Timestamp: timestamp,
		Signature: sig,
	}

	// Make the POST request to the Allowlist API
	resp, err := request.GetLivelinessFileAPI(reqBody, h.Config)
	if err != nil {
		http.Error(w, "API call failed, "+err.Error(), http.StatusInternalServerError)
		return
	}

	payloadBody, err := io.ReadAll(&resp.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("\n\nAPI response: %s", payloadBody)

	if resp.StatusCode != 200 {
		var rw response.ResponseWrapper
		err = rw.FromJSON(payloadBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("\n\nPayload error response: %v", rw)
		// Write the API response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write(payloadBody)
		return
		// If you want to store the error response somewhere, do that here using the wrapper response.
		// e.g., {
		// 	"code": 500,
		// 	"message": "There was an error processing your request. It has been logged (ID f4f981f01de5a056)."
		// }
		// Where, rw.code = 500 and rw.message = "the message seen above"
	}

	// Response 200 OK, will be binary data, as in the blob file.
	// Check response content-type to assign correct file extension.
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type not found in the response", http.StatusBadRequest)
		return
	}

	mimeType := strings.Split(contentType, ";")[0]
	fileExtension := ".mp4"
	if mimeType == "image/gif" {
		fileExtension = ".gif"
	}
	fmt.Printf("The file extension for MIME type '%s' is '%s'\n", mimeType, fileExtension)

	// Set file name, i.e. txID_gifResultID_liveliness.png
	fileName := fmt.Sprintf("%s_%s_liveliness%s", f.TxID, f.ResultGIFID, fileExtension)
	err = os.WriteFile(fileName, payloadBody, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("File saved as '%s'\n", fileName)

	// Save the binary data to a file (optional)
	err = os.WriteFile(fileName, payloadBody, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write file to webapp
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(payloadBody)
}
