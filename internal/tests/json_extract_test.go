package tests

import (
	"encoding/json"
	"testing"
)

func extractToken(t *testing.T, body []byte) string {
	t.Helper()
	var resp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &resp)
	return resp.Data.AccessToken
}

func extractID(t *testing.T, body []byte) int {
	t.Helper()
	var resp struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &resp)
	return resp.Data.ID
}