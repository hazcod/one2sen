package onepassword

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	onePasswordTimestampFormat = "2006-01-02T15:04:05-07:00"
)

func toJson(obj interface{}) (string, error) {
	b, err := json.Marshal(&obj)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ConvertUsageToMap(_ *logrus.Logger, items []Item) ([]map[string]string, error) {
	logs := make([]map[string]string, len(items))

	var err error

	for i, item := range items {
		cols := make(map[string]string)

		cols["LogType"] = "Usage"

		cols["TimeGenerated"] = item.Timestamp.UTC().Format(time.RFC3339)
		//cols["UUID"] = item.UUID

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
			"Action":    item.Action,
			"VaultUUID": item.VaultUUID,
			"ItemUUID":  item.ItemUUID,
		}
		cols["Data"], err = toJson(custom)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal data: %v", err)
		}

		logs[i] = cols
	}

	return logs, err
}

func ConvertEventToMap(_ *logrus.Logger, events []Event) ([]map[string]string, error) {
	logs := make([]map[string]string, len(events))

	var err error

	for i, event := range events {
		cols := make(map[string]string)

		cols["LogType"] = "Event"

		cols["TimeGenerated"] = event.Timestamp.UTC().Format(time.RFC3339)
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

		cols["TimeGenerated"] = event.Timestamp.UTC().Format(time.RFC3339)
		//cols["UUID"] = event.UUID

		cols["User"], err = toJson(event.ActorUUID)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal User: %v", err)
		}

		cols["Location"], err = toJson(event.Location)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal Location: %v", err)
		}

		// specific columns
		custom := map[string]string{
			"Action":      event.Action,
			"ActorUUID":   event.ActorUUID,
			"ObjectType":  event.ObjectType,
			"ObjectUUID":  event.ObjectUUID,
			"SessionUUID": event.Session.UUID,
		}

		cols["Data"], err = toJson(custom)
		if err != nil {
			return nil, fmt.Errorf("could not json marshal data: %v", err)
		}

		logs[i] = cols
	}

	return logs, err
}
