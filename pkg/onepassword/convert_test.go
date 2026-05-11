package onepassword

import "testing"

func TestConvertUsageToMapParsesOffsetTimestamps(t *testing.T) {
	items := []Item{
		{
			Timestamp: "2026-05-11T11:46:01+02:00",
			User: User{
				UUID:  "user-1",
				Name:  "Jane Doe",
				Email: "jane@example.com",
			},
			Location: Location{
				City:    "Brussels",
				Country: "Belgium",
			},
		},
	}
	logs, err := ConvertUsageToMap(nil, items)
	if err != nil {
		t.Fatalf("ConvertUsageToMap returned error: %v", err)
	}
	if got, want := logs[0]["TimeGenerated"], "2026-05-11T09:46:01Z"; got != want {
		t.Fatalf("TimeGenerated = %q, want %q", got, want)
	}
}
