/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package request

import (
	"encoding/json"
	"fmt"

	client "github.com/caelisco/http-client"
	"github.com/intergreatme/remote-kyc-util/config"
)

type RequestPayload struct {
	Payload   string `json:"payload"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// ToJSON serializes the RequestPayload struct to JSON
func (r RequestPayload) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes JSON data into the RequestPayload struct
func (r *RequestPayload) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}

func AllowlistAPI(payload RequestPayload, cnf config.Configuration) (client.Response, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return client.Response{}, err
	}

	uri := fmt.Sprintf("%sv2/allowlist/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "https://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
		return client.Response{}, err
	}

	return resp, nil
}

func UpdateAPI(payload RequestPayload, cnf config.Configuration) (client.Response, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return client.Response{}, err
	}

	uri := fmt.Sprintf("%supdate/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "https://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
		return client.Response{}, err
	}

	return resp, nil
}


func CancelAPI(payload RequestPayload, cnf config.Configuration) (client.Response, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return client.Response{}, err
	}

	uri := fmt.Sprintf("%scancel/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "https://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
		return client.Response{}, err
	}

	return resp, nil
}


func GetFileAPI(payload RequestPayload, cnf config.Configuration) (client.Response, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return client.Response{}, err
	}

	uri := fmt.Sprintf("%sgetfile/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "http://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
		return client.Response{}, err
	}

	return resp, nil
}

func GetLivelinessFileAPI(payload RequestPayload, cnf config.Configuration) (client.Response, error) {
	// Compress the JSON payload with gzip
	out, err := payload.ToJSON()
	if err != nil {
		return client.Response{}, err
	}

	uri := fmt.Sprintf("%sgetlivelinessfile/%s", cnf.URL, cnf.CompanyID)

	opt := client.RequestOptions{
		Compression:    client.CompressionGzip,
		ProtocolScheme: "http://",
	}

	opt.AddHeader("Content-Type", "application/json")

	resp, err := client.Post(uri, out, opt)
	if err != nil {
		fmt.Println(err)
		return client.Response{}, err
	}

	return resp, nil
}
