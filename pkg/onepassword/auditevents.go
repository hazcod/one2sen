package onepassword

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type auditEventResponse struct {
	Cursor  string       `json:"cursor"`
	HasMore bool         `json:"has_more"`
	Items   []AuditEvent `json:"items"`
}

type Session struct {
	UUID       string `json:"uuid"`
	LoginTime  string `json:"login_time"`
	DeviceUUID string `json:"device_uuid"`
	IP         string `json:"ip"`
}

type AuditEvent struct {
	UUID       string    `json:"uuid"`
	Timestamp  time.Time `json:"timestamp"`
	ActorUUID  string    `json:"actor_uuid"`
	Action     string    `json:"action"`
	ObjectType string    `json:"object_type"`
	ObjectUUID string    `json:"object_uuid"`
	AuxID      int       `json:"aux_id"`
	AuxUUID    string    `json:"aux_uuid"`
	AuxInfo    string    `json:"aux_info"`
	Session    Session   `json:"session"`
	Location   Location  `json:"location"`
}

func (p *OnePassword) GetAuditEvents(lookBackDays uint) ([]AuditEvent, error) {
	items := make([]AuditEvent, 0)

	startTime := time.Now().UTC().AddDate(0, 0, -1*int(lookBackDays))
	endTime := time.Now().UTC()

	hasMore := true
	cursor := ""

	for hasMore {
		p.Logger.Debug("fetching audit events")

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

		usagesRequest, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/auditevents", eventsURL), bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, fmt.Errorf("could not create usage request: %v", err)
		}

		usagesRequest.Header.Set("Content-Type", "application/json")
		usagesRequest.Header.Set("Authorization", "Bearer "+p.apiToken)

		usagesResponse, usagesError := p.httpClient.Do(usagesRequest)
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

		var resp auditEventResponse

		if err := json.Unmarshal(usagesBody, &resp); err != nil {
			return nil, fmt.Errorf("could not decode usage response: %v", err)
		}

		hasMore = resp.HasMore
		cursor = resp.Cursor

		items = append(items, resp.Items...)
	}

	p.Logger.WithField("total", len(items)).Debug("retrieved audit events")

	return items, nil
}
