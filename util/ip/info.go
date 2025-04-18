package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Info struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

// Get public IP and information about it.
func Get() (Info, error) {
	resp, err := http.Get("https://ipinfo.io/") // IPinfo returns IPv4 or IPv6 based on the protocol stack from your request source.
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
