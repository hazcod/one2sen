package onepassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type eventResponse struct {
	Cursor  string  `json:"cursor"`
	HasMore bool    `json:"has_more"`
	Items   []Event `json:"items"`
	Error   struct {
		Message string `json:"Message"`
	} `json:"Error"`
}

type TargetUser struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Event struct {
	UUID        string      `json:"uuid"`
	SessionUUID string      `json:"session_uuid"`
	Timestamp   time.Time   `json:"timestamp"`
	Country     string      `json:"country"`
	Category    string      `json:"category"`
	Type        string      `json:"type"`
	Details     interface{} `json:"details"`
	TargetUser  TargetUser  `json:"target_user"`
	Client      Client      `json:"client"`
	Location    Location    `json:"location"`
}

func (e *Event) IsOK() bool {
	return strings.Contains(strings.ToLower(e.Type), "_ok")
}

func (p *OnePassword) GetSigninEvents(lookBackDays uint) ([]Event, error) {
	items := make([]Event, 0)

	startTime := time.Now().UTC().AddDate(0, 0, -1*int(lookBackDays))
	endTime := time.Now().UTC()

	hasMore := true
	cursor := ""

	for hasMore {
		p.Logger.Debug("fetching signin events")

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

		signinRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/signinattempts", eventsURL), bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("could not create signin request: %v", err)
		}

		signinRequest.Header.Set("Content-Type", "application/json")
		signinRequest.Header.Set("Authorization", "Bearer "+p.apiToken)

		signinResponse, err := httpClient.Do(signinRequest)
		if err != nil {
			return nil, fmt.Errorf("could not fetch signins: %v", err)
		}
		defer signinResponse.Body.Close()

		if signinResponse.StatusCode > 399 {
			return nil, fmt.Errorf("returned status code: %d", signinResponse.StatusCode)
		}

		signinsBody, err := io.ReadAll(signinResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read signin response body: %v", err)
		}

		p.Logger.Tracef("%s", signinsBody)

		var resp eventResponse

		if err := json.Unmarshal(signinsBody, &resp); err != nil {
			return nil, fmt.Errorf("could not decode usage response: %v", err)
		}

		if resp.Error.Message != "" {
			return nil, fmt.Errorf("returned error: %v", resp.Error.Message)
		}

		hasMore = resp.HasMore
		cursor = resp.Cursor

		items = append(items, resp.Items...)
	}

	p.Logger.WithField("total", len(items)).Debug("retrieved signin events")

	return items, nil
}
