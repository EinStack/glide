package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// Payload represents the payload data structure
type Payload map[string]interface{}

// RenamePayloadKeys renames payload keys based on the specified mapping. If a key is not present in the
// mapping, the key and its value will remain unchanged.
// TODO: This functionality present in buildAPIConfig.go. Bring over?
func RenamePayloadKeys(payload Payload, mapping map[string]string) Payload {
	newPayload := make(Payload)
	for k, v := range payload {
		newKey, ok := mapping[k]
		if !ok {
			newKey = k
		}
		newPayload[newKey] = v
	}
	return newPayload
}

// SendRequest sends an HTTP request to a specific URL path with given headers and payload.
func SendRequest(headers map[string]string, baseURL string, path string, payload Payload) (map[string]interface{}, error) {
	client := &http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(result["error"].(string))
	}

	return result, nil
}
