package gohassapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type HassClient struct {
	Endpoint string
	Token    string

	client *http.Client
}

func NewClient(endpoint string, token string) *HassClient {
	return &HassClient{
		Endpoint: endpoint,
		Token:    token,
		client:   &http.Client{},
	}
}

func (hc *HassClient) get(path string) ([]byte, error) {
	url, err := url.JoinPath(hc.Endpoint, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil /* body */)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", hc.Token))
	req.Header.Add("content-type", "application/json")
	resp, err := hc.client.Do(req)

	switch resp.StatusCode {
	case 200:
	case 201:
	case 400:
		return nil, fmt.Errorf("API returned HTTP 400 (Bad Request)")
	case 401:
		return nil, fmt.Errorf("API returned HTTP 401 (Unauthorized)")
	case 404:
		return nil, fmt.Errorf("API returned HTTP 404 (Not Found)")
	case 405:
		return nil, fmt.Errorf("API returned HTTP 405 (Method Not Allowed)")
	default:
		return nil, fmt.Errorf("API returned an unknown HTTP error code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}

func genericGet[T any](hc *HassClient, path string) (T, error) {
	resp, err := hc.get(path)
	if err != nil {
		return *new(T), err
	}
	var r T
	if err := json.Unmarshal(resp, &r); err != nil {
		slog.Debug("Failed to unmarshal JSON object", err.Error(), resp)
		return *new(T), err
	}
	return r, nil

}

// Queries the /api/ endpoint and reports if the API is up and
// running.
func (hc *HassClient) Check() (string, error) {
	ac, err := genericGet[Check](hc, "/")
	if err != nil {
		return "", err
	}
	return ac.Message, nil
}

// Queries the /states endpoint and returns a list of entities with their associated states
func (hc *HassClient) States() ([]State, error) {
	return genericGet[[]State](hc, "/states")
}

// Queries the /services endpoint and returns a list of service objects.
func (hc *HassClient) Services() ([]ServiceDomain, error) {
	return genericGet[[]ServiceDomain](hc, "/services")
}
