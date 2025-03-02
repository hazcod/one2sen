package onepassword

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	onePasswordEventTimestampFormat = "2006-01-02T15:04:05.99999999Z"
	onePasswordTimestampFormat      = "2006-01-02T15:04:05-07:00"
	iso8601Format                   = "2006-01-02T15:04:05Z"
)

func toJson(obj interface{}) (string, error) {
	switch _, ok := obj.(string); {
	case ok:
		return obj.(string), nil
	}

	b, err := json.Marshal(&obj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal object: %w", err)
	}

	return string(b), nil
}

func ConvertUsageToMap(_ *logrus.Logger, items []Item) ([]map[string]string, error) {
	logs := make([]map[string]string, len(items))

	var err error

	for i, item := range items {
		cols := make(map[string]string)

		cols["LogType"] = "Usage"

		// YYYY-MM-DDThh:mm:ssZ
		timeGenerated, err := time.Parse(onePasswordEventTimestampFormat, item.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 1P event timestamp: %w", err)
		}

		cols["TimeGenerated"] = timeGenerated.UTC().Format(iso8601Format)

		cols["User"], err = toJson(item.User)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal User: %v", err)
		}

		cols["Client"], err = toJson(item.Client)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Client: %v", err)
		}

		cols["Location"], err = toJson(item.Location)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Location: %v", err)
		}

		// specific columns
		custom := map[string]string{
			"Action":     item.Action,
			"VaultUUID":  item.VaultUUID,
			"ItemUUID":   item.ItemUUID,
			"ActorUUID":  item.User.UUID,
			"ActorName":  item.User.Name,
			"ActorEmail": item.User.Email,
			"City":       item.Location.City,
			"Country":    item.Location.Country,
		}
		cols["Data"], err = toJson(custom)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal data: %v", err)
		}

		logs[i] = cols
	}

	return logs, err
}

func ConvertSigninToMap(_ *logrus.Logger, events []Event) ([]map[string]string, error) {
	logs := make([]map[string]string, len(events))

	var err error

	for i, event := range events {
		cols := make(map[string]string)

		cols["LogType"] = "Event"

		timeGenerated, err := time.Parse(onePasswordEventTimestampFormat, event.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 1P event timestamp: %w", err)
		}

		cols["TimeGenerated"] = timeGenerated.UTC().Format(iso8601Format)
		//cols["UUID"] = event.UUID

		cols["User"], err = toJson(event.TargetUser)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal User: %v", err)
		}

		cols["Client"], err = toJson(event.Client)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Client: %v", err)
		}

		cols["Location"], err = toJson(event.Location)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Location: %v", err)
		}

		// specific columns
		eventDetails, err := toJson(event.Details)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Details: %v", err)
		}
		custom := map[string]string{
			"OK":          fmt.Sprintf("%t", event.IsOK()),
			"Details":     eventDetails,
			"SessionUUID": event.SessionUUID,
			"EventType":   event.Type,
			"ActorUUID":   event.TargetUser.UUID,
			"ActorName":   event.TargetUser.Name,
			"ActorEmail":  event.TargetUser.Email,
			"City":        event.Location.City,
			"Country":     event.Location.Country,
		}

		cols["Data"], err = toJson(custom)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal data: %v", err)
		}

		logs[i] = cols
	}

	return logs, err
}

func ConvertAuditEventToMap(_ *logrus.Logger, audits []AuditEvent) ([]map[string]string, error) {
	logs := make([]map[string]string, len(audits))

	var err error

	for i, event := range audits {
		cols := make(map[string]string)

		cols["LogType"] = "Audit"

		timeGenerated, err := time.Parse(onePasswordEventTimestampFormat, event.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse 1P event timestamp: %w", err)
		}

		cols["TimeGenerated"] = timeGenerated.UTC().Format(iso8601Format)
		//cols["UUID"] = event.UUID

		cols["User"], err = toJson(event.ActorUUID)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal User: %v", err)
		}

		cols["Location"], err = toJson(event.Location)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Location: %v", err)
		}

		auxDetails, err := toJson(event.AuxDetails)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal auxDetails: %v", err)
		}

		// specific columns
		custom := map[string]string{
			"Action":      event.Action,
			"ActorUUID":   event.ActorUUID,
			"ActorName":   event.ActorDetails.Name,
			"ActorEmail":  event.ActorDetails.Email,
			"City":        event.Location.City,
			"Country":     event.Location.Country,
			"ObjectType":  event.ObjectType,
			"ObjectUUID":  event.ObjectUUID,
			"SessionUUID": event.Session.UUID,
			"AuxUUID":     event.AuxUUID,
			"AuxDetails":  auxDetails,
		}

		cols["Data"], err = toJson(custom)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal data: %v", err)
		}

		logs[i] = cols
	}

	return logs, err
}
