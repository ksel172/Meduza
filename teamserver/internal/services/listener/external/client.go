package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

// ListenerClient handles communication with external listener implementations
type ListenerClient struct {
	apiKey      string
	httpClient  *http.Client
	listenerDal dal.IListenerDAL
}

func NewListenerClient(apiKey string, listenerDal dal.IListenerDAL) *ListenerClient {
	return &ListenerClient{
		apiKey:      apiKey,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		listenerDal: listenerDal,
	}
}

// AddListener contacts the external listener controller to add a new listener
func (lc *ListenerClient) AddListener(listener models.Listener, url string) error {
	data := map[string]string{
		"id":     listener.ID,
		"config": string(listener.RawConfig),
	}

	return lc.sendRequest(http.MethodPost, fmt.Sprintf("%s/listeners", url), data)
}

// UpdateListener contacts the external listener controller to update a listener
func (lc *ListenerClient) UpdateListener(listenerID string, config string, url string) error {
	data := map[string]string{
		"id":     listenerID,
		"config": config,
	}

	return lc.sendRequest(http.MethodPut, fmt.Sprintf("%s/listeners/%s", url, listenerID), data)
}

// StartListener contacts the external listener controller to start a listener
func (lc *ListenerClient) StartListener(listenerID string, url string) error {
	data := map[string]string{
		"id": listenerID,
	}

	return lc.sendRequest(http.MethodPost, fmt.Sprintf("%s/listeners/%s/start", url, listenerID), data)
}

// StopListener contacts the external listener controller to stop a listener
func (lc *ListenerClient) StopListener(listenerID string, url string) error {
	data := map[string]string{
		"id": listenerID,
	}

	return lc.sendRequest(http.MethodPost, fmt.Sprintf("%s/listeners/%s/stop", url, listenerID), data)
}

// TerminateListener contacts the external listener controller to terminate listener
func (lc *ListenerClient) TerminateListener(listenerID string, url string) error {
	data := map[string]string{
		"id": listenerID,
	}

	return lc.sendRequest(http.MethodDelete, fmt.Sprintf("%s/listeners/%s", url, listenerID), data)
}

// Helper function to send requests
func (lc *ListenerClient) sendRequest(method, url string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))
	// We should maybe add a UserAgent .env var and set it here

	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	// Optionally parse response body if needed
	// var response struct {
	//     Status string `json:"status"`
	//     Message string `json:"message"`
	// }
	// if err := json.Unmarshal(body, &response); err != nil {
	//     return fmt.Errorf("failed to parse response: %w", err)
	// }

	return nil
}

// GetListenerStatus queries the external controller for the current status of a listener
func (lc *ListenerClient) GetListenerStatus(listenerID string, url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqURL := fmt.Sprintf("%s/listeners/%s/status", url, listenerID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to prepare request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lc.apiKey))

	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	// Parse response
	var response struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Status, nil
}

// Maybe it's worth implementing a health check from server to controller instead of the other way around

// func (lc *ListenerClient) HealthCheck(url string) error {

// }
