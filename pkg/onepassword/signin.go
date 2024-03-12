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
	Timestamp   string      `json:"timestamp"`
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

	now := time.Now().UTC()
	startTime := now.AddDate(0, 0, -1*int(lookBackDays))

	round := 0
	hasMore := true
	cursor := ""

	for hasMore {
		round += 1
		p.Logger.WithField("round", round).Debug("fetching signin events")

		payload := eventRequest{}
		if cursor != "" {
			payload.Cursor = cursor
		} else {
			payload.Limit = maxFetch
			payload.StartTime = startTime.Format(onePasswordTimestampFormat)
			payload.EndTime = now.Format(onePasswordTimestampFormat)
		}

		payloadBytes, err := json.Marshal(&payload)
		if err != nil {
			return nil, fmt.Errorf("could not encode payload: %v", err)
		}

		signinRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/signinattempts", p.apiURL), bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("could not create signin request: %v", err)
		}

		signinRequest.Header.Set("Content-Type", "application/json")
		signinRequest.Header.Set("Authorization", "Bearer "+p.apiToken)

		signinResponse, err := p.httpClient.Do(signinRequest)
		if err != nil {
			return nil, fmt.Errorf("could not fetch signins: %v", err)
		}

		if signinResponse.StatusCode > 399 {
			_ = signinResponse.Body.Close()
			return nil, fmt.Errorf("returned status code: %d", signinResponse.StatusCode)
		}

		signinsBody, err := io.ReadAll(signinResponse.Body)
		_ = signinResponse.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("could not read signin response body: %v", err)
		}

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
