/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */
package cancel

import "encoding/json"

type Cancel struct {
	TxID    string `json:"tx_id"`
	Comment string `json:"comment,omitempty"`
}

func (a Cancel) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Cancel) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}
