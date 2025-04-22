package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Info generated from https://ip-api.com/docs/api:json
type Info struct {
	IP       string  `json:"query"`
	Country  string  `json:"country"`
	Region   string  `json:"regionName"`
	City     string  `json:"city"`
	District string  `json:"district"`
	Zip      string  `json:"zip"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`
	Isp      string  `json:"isp"`
	Org      string  `json:"org"`
	As       string  `json:"as"`
}

// Get public IP.
func Get() (string, error) {
	resp, err := http.Get("https://icanhazip.com/")
	if err != nil {
		return "", fmt.Errorf("failed to request api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return strings.TrimSpace(string(body)), nil
}

// Query IP information.
func Query(ip string) (Info, error) {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s?lang=en&fields=%s", ip, "536569"))
	if err != nil {
		return Info{}, fmt.Errorf("failed to request api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Info{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Info{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var info Info
	err = json.Unmarshal(body, &info)
	if err != nil {
		return Info{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return info, nil
}
