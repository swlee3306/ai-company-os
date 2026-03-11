package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type openclawInvokeRequest struct {
	Tool       string         `json:"tool"`
	Action     string         `json:"action,omitempty"`
	Args       map[string]any `json:"args,omitempty"`
	SessionKey string         `json:"sessionKey,omitempty"`
}

type openclawInvokeResponse struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func openclawInvoke(tool string, args map[string]any) (json.RawMessage, error) {
	url := os.Getenv("OPENCLAW_GATEWAY_URL")
	if url == "" {
		url = "http://127.0.0.1:18789"
	}
	token := os.Getenv("OPENCLAW_GATEWAY_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("OPENCLAW_GATEWAY_TOKEN is required for openclaw runner")
	}
	reqBody := openclawInvokeRequest{Tool: tool, Args: args}
	b, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", url+"/tools/invoke", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("openclaw invoke failed: %s", string(body))
	}
	var out openclawInvokeResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if !out.OK {
		if out.Error != nil {
			return nil, fmt.Errorf("openclaw invoke error: %s", out.Error.Message)
		}
		return nil, fmt.Errorf("openclaw invoke error")
	}
	return out.Result, nil
}
