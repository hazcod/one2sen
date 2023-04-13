package onepassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type usageResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Items   []Item `json:"items"`
	Error   struct {
		Message string `json:"Message"`
	} `json:"Error"`
}

type User struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Client struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	PlatformName    string `json:"platform_name"`
	PlatformVersion string `json:"platform_version"`
	OsName          string `json:"os_name"`
	OsVersion       string `json:"os_version"`
	IPAddress       string `json:"ip_address"`
}

type Location struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Item struct {
	UUID        string    `json:"uuid"`
	Timestamp   time.Time `json:"timestamp"`
	UsedVersion int       `json:"used_version"`
	VaultUUID   string    `json:"vault_uuid"`
	ItemUUID    string    `json:"item_uuid"`
	User        User      `json:"user"`
	Client      Client    `json:"client"`
	Location    Location  `json:"location"`
	Action      string    `json:"action"`
}

func (p *OnePassword) GetUsage(lookBackDays uint) ([]Item, error) {
	items := make([]Item, 0)

	startTime := time.Now().UTC().AddDate(0, 0, -1*int(lookBackDays))
	endTime := time.Now().UTC()

	hasMore := true
	cursor := ""

	for hasMore {
		p.Logger.Debug("fetching usage events")

		payload := eventRequest{
			Limit:     maxFetch,
			Cursor:    cursor,
			StartTime: startTime.Format(time.RFC3339),
			EndTime:   endTime.Format(time.RFC3339),
		}

		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, fmt.Errorf("could not encode payload: %v", err)
		}

		p.Logger.Debugf("%s", payloadBytes)

		usagesRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/itemusages", eventsURL), bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("could not create usage request: %v", err)
		}

		usagesRequest.Header.Set("Content-Type", "application/json")
		usagesRequest.Header.Set("Authorization", "Bearer "+p.apiToken)

		usagesResponse, usagesError := httpClient.Do(usagesRequest)
		if usagesError != nil {
			return nil, fmt.Errorf("could not fetch usage: %v", err)
		}

		defer usagesResponse.Body.Close()

		if usagesResponse.StatusCode > 399 {
			return nil, fmt.Errorf("returned status code: %d", usagesResponse.StatusCode)
		}

		usagesBody, err := io.ReadAll(usagesResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read usage: %v", err)
		}

		var resp usageResponse

		if err := json.Unmarshal(usagesBody, &resp); err != nil {
			return nil, fmt.Errorf("could not decode usage response: %v", err)
		}

		p.Logger.Tracef("%+v", resp)

		if resp.Error.Message != "" {
			return nil, fmt.Errorf("returned error: %v", resp.Error.Message)
		}

		hasMore = resp.HasMore
		cursor = resp.Cursor

		items = append(items, resp.Items...)
	}

	p.Logger.WithField("total", len(items)).Debug("retrieved usage events")

	return items, nil
}
