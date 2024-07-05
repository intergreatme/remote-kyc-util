package file

// Payload represents the structure of the payload
type FilePayload struct {
	TxID         string `json:"tx_id"`
	DocumentType string `json:"document_type"`
	FileType     string `json:"file_type"`
}

type LivelinessPayload struct {
	TxID        string `json:"tx_id"`
	ResultGIFID string `json:"result_gif_id"`
}
